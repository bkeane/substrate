# https://docs.aws.amazon.com/lambda/latest/dg/images-create.html
resource "aws_ecr_repository_policy" "cross_account_access" {
    for_each = local.is_hub ? local.artifacts : toset([])
    repository = each.value
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Sid = "CrossAccountPermission"
                Effect = "Allow"
                Principal = {
                    AWS = concat([
                        for spoke in var.spoke_accounts:
                            "arn:aws:iam::${spoke.account_id}:root"
                    ])
                }
                Action = [
                    "ecr:BatchGetImage",
                    "ecr:GetDownloadUrlForLayer"
                ]
            },
            {
                Sid = "LambdaECRImageRetrievalPolicy"
                Effect = "Allow"
                Action = [
                    "ecr:BatchGetImage",
                    "ecr:GetDownloadUrlForLayer"
                ]
                Principal = {
                    Service = "lambda.amazonaws.com"
                }
                Condition = {
                    StringLike = {
                        "aws:sourceARN": concat([
                            for spoke in var.spoke_accounts: 
                                "arn:aws:lambda:${data.aws_region.current.name}:${spoke.account_id}:function:*"
                        ])
                    }
                }
            }
        ]
    })
}

