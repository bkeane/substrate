
#
# GitHub OIDC Provider
#

resource "aws_iam_openid_connect_provider" "github" {
  count = var.create_oidc_provider ? 1 : 0
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  # Hex-encoded SHA-1 hash of the X.509 domain certificate
  # https://github.blog/changelog/2023-06-27-github-actions-update-on-oidc-integration-with-aws/
  thumbprint_list = [
        "6938fd4d98bab03faadb97b34396831e3780aea1",
        "1c58a3a8518e8759bf075b76b750d4f2df264fcd"
    ]
}

data "aws_iam_openid_connect_provider" "github" {
  count = var.create_oidc_provider ? 0 : 1
  arn = "arn:aws:iam::${var.hub_account.account_id}:oidc-provider/token.actions.githubusercontent.com"
}

#
# GitHub Actions Role
#

resource "aws_iam_role" "github" {
  name                  = "${var.name}-oidc-role"
  description           = "used by github actions"
  assume_role_policy    = data.aws_iam_policy_document.trust.json
  force_detach_policies = true
}

data "aws_iam_policy_document" "trust" {
  statement {
    principals {
      type        = "Federated"
      identifiers = var.create_oidc_provider ? [aws_iam_openid_connect_provider.github[0].arn] : [data.aws_iam_openid_connect_provider.github[0].arn]
    }

    actions = ["sts:AssumeRoleWithWebIdentity"]

    condition {
        test     = "StringEquals"
        variable = "token.actions.githubusercontent.com:aud"
        values   = ["sts.amazonaws.com"]
    }

    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values = [
        for org in local.organizations:
            "repo:${org}/*"
      ]
    }
  }
}

resource "aws_iam_role_policy_attachment" "github" {
  role       = aws_iam_role.github.name
  policy_arn = aws_iam_policy.github.arn
}

#
# GitHub Actions Policy
#

resource "aws_iam_policy" "github" {
  name        = "${var.name}-oidc-policy"
  description = "used by github actions"
  policy      = data.aws_iam_policy_document.github.json
}


data "aws_iam_policy_document" "github" {
  statement {
    sid = "AllowEcrLogin"
    effect = "Allow"
    resources = ["*"]
    actions   = ["ecr:GetAuthorizationToken"]
  }

  statement {
    sid = "AllowEcrRepositoryAccess"
    effect = "Allow"
    actions = [
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "ecr:BatchCheckLayerAvailability",
      "ecr:PutImage",
      "ecr:InitiateLayerUpload",
      "ecr:UploadLayerPart",
      "ecr:CompleteLayerUpload",
      # "ecr:BatchDeleteImage",
    ]
    resources = [
        for image_path in local.image_paths :
            "arn:aws:ecr:${var.hub_account.account_region}:${var.hub_account.account_id}:repository/${image_path}"
    ]
  }

  statement {
    sid = "AllowEventBridgePutEvents"
    effect = "Allow"
    actions = ["events:PutEvents"]
    resources = [
      for bus_name in var.bus_names:
        "arn:aws:events:${var.hub_account.account_region}:${var.hub_account.account_id}:eventbus/${bus_name}"
    ]
  }
}
