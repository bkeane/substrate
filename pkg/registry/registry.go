package registry

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type MediaType map[string]map[string]string

const (
	DOCKER_MANIFEST_INDEX = "application/vnd.docker.distribution.manifest.list.v2+json"
	DOCKER_MANIFEST       = "application/vnd.docker.distribution.manifest.v2+json"
	OCI_MANIFEST_INDEX    = "application/vnd.oci.image.index.v1+json"
	OCI_MANIFEST          = "application/vnd.oci.image.manifest.v1+json"
	OCI_CONFIG_MANIFEST   = "application/vnd.oci.image.config.v1+json"
)

type Registry struct {
	url   string
	token string
}

type Catalogue struct {
	Repositories []string `json:"repositories"`
}

type Tags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ImageIndex struct {
	Manifests []struct {
		MediaType string             `json:"mediaType"`
		Digest    string             `json:"digest"`
		Size      int                `json:"size"`
		Platform  ImageIndexPlatform `json:"platform"`
	} `json:"manifests"`
}

type ImageIndexPlatform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

type ImageManifest struct {
	Config struct {
		Digest string `json:"digest"`
	} `json:"config"`
}

type ImageConfig struct {
	Config struct {
		Labels map[string]string `json:"labels"`
	} `json:"config"`
	Architecture string `json:"architecture"`
}

type ImagePointer struct {
	ImageConfig
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Digest     string `json:"digest"`
	Uri        string `json:"uri"`
}

func GetImageFromUri(ctx context.Context, awsconfig aws.Config, imageUri string) (ImagePointer, error) {
	if !strings.Contains(imageUri, "@") {
		return ImagePointer{}, fmt.Errorf("invalid image uri: %s", imageUri)
	}

	if !strings.HasPrefix(imageUri, "http") {
		imageUri = fmt.Sprintf("https://%s", imageUri)
	}

	uri, digest := strings.Split(imageUri, "@")[0], strings.Split(imageUri, "@")[1]
	url, err := url.Parse(uri)
	if err != nil {
		return ImagePointer{}, err
	}

	reg, err := Init(ctx, awsconfig, url.Host)
	if err != nil {
		return ImagePointer{}, err
	}

	return reg.GetImageFromName(ctx, strings.TrimPrefix(url.Path, "/"), digest)
}

func Init(ctx context.Context, awsconfig aws.Config, url string) (*Registry, error) {
	r := &Registry{}
	r.url = url

	ecrc := ecr.NewFromConfig(awsconfig)
	ecrauth, err := ecrc.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, err
	}

	r.token = *ecrauth.AuthorizationData[0].AuthorizationToken

	return r, nil
}

func (r *Registry) GetRepositories(ctx context.Context) (Catalogue, error) {
	url := fmt.Sprintf("https://%s/v2/_catalog", r.url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+r.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Catalogue{}, err
	}
	defer resp.Body.Close()

	var catalogue Catalogue
	if err := json.NewDecoder(resp.Body).Decode(&catalogue); err != nil {
		return Catalogue{}, err
	}

	return catalogue, nil
}

func (r *Registry) GetTags(ctx context.Context, repository string) (Tags, error) {
	url := fmt.Sprintf("https://%s/v2/%s/tags/list", r.url, repository)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+r.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Tags{}, err
	}
	defer resp.Body.Close()

	var tags Tags
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return Tags{}, err
	}

	return tags, nil
}

func (r *Registry) GetImageFromName(ctx context.Context, repository string, reference string) (ImagePointer, error) {
	jsonString, err := r.DigImage(ctx, repository, reference)
	if err != nil {
		return ImagePointer{}, err
	}

	var pointer ImagePointer
	if err := json.Unmarshal([]byte(jsonString), &pointer); err != nil {
		return ImagePointer{}, err
	}

	return pointer, nil
}

func (r *Registry) DigImage(ctx context.Context, repository string, reference string) (string, error) {
	resp, err := r.GetManifest(ctx, repository, reference)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.Header.Get("Content-Type") {
	case DOCKER_MANIFEST_INDEX, OCI_MANIFEST_INDEX:
		var index ImageIndex
		if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
			return "", err
		}

		var digest string
		// Set default architecture to that which exists
		for _, manifest := range index.Manifests {
			if manifest.Platform.Architecture != "unknown" {
				digest = manifest.Digest
			}
		}

		// set default architecture to arm64 if exists
		for _, manifest := range index.Manifests {
			if manifest.Platform.Architecture == "arm64" {
				digest = manifest.Digest
				break
			}
		}

		return r.DigImage(ctx, repository, digest)

	case DOCKER_MANIFEST, OCI_MANIFEST:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		hasher := sha256.New()
		if _, err := io.Copy(hasher, bytes.NewReader(bodyBytes)); err != nil {
			return "", err
		}
		digest := "sha256:" + hex.EncodeToString(hasher.Sum(nil))

		var manifest ImageManifest
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&manifest); err != nil {
			return "", err
		}

		resp, err := r.GetConfig(ctx, repository, manifest.Config.Digest)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var pointer ImagePointer
		if err := json.NewDecoder(resp.Body).Decode(&pointer); err != nil {
			return "", err
		}

		pointer.Registry = r.url
		pointer.Repository = repository
		pointer.Digest = digest
		pointer.Uri = fmt.Sprintf("%s/%s@%s", r.url, repository, digest)

		jsonBytes, err := json.Marshal(pointer)
		if err != nil {
			return "", err
		}

		return string(jsonBytes), nil

	default:
		return "", fmt.Errorf("unknown content type %s", resp.Header.Get("Content-Type"))
	}
}

func (r *Registry) GetManifest(ctx context.Context, repository string, reference string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", r.url, repository, reference)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+r.token)
	req.Header.Add("Accept", strings.Join([]string{DOCKER_MANIFEST_INDEX, DOCKER_MANIFEST, OCI_MANIFEST_INDEX, OCI_MANIFEST}, ","))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get image manifest: status %d", resp.StatusCode)
	}

	return resp, nil
}

func (r *Registry) GetConfig(ctx context.Context, repository string, reference string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/v2/%s/blobs/%s", r.url, repository, reference)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+r.token)
	req.Header.Add("Accept", strings.Join([]string{OCI_CONFIG_MANIFEST}, ","))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get image config: status %d", resp.StatusCode)
	}

	return resp, nil
}
