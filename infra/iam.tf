data "aws_iam_policy_document" "ec2_assume" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "ec2" {
  name               = "app-ec2-role"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume.json
}

# Allow SSM Session Manager and Parameter Store read
resource "aws_iam_role_policy_attachment" "ssm_core" {
  role       = aws_iam_role.ec2.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

# Read the specific parameter and decrypt with AWS managed key for SSM
data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_iam_policy" "ssm_read" {
  name   = "app-ssm-read"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = ["ssm:GetParameter", "ssm:GetParameters"],
        Resource = [aws_ssm_parameter.db_password.arn]
      },
      {
        Effect = "Allow",
        Action = ["kms:Decrypt"],
        Resource = [
          # AWS managed key alias for SSM in this account/region
          "arn:${data.aws_partition.current.partition}:kms:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:alias/aws/ssm"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ssm_read" {
  role       = aws_iam_role.ec2.name
  policy_arn = aws_iam_policy.ssm_read.arn
}

resource "aws_iam_instance_profile" "ec2" {
  name = "app-ec2-instance-profile"
  role = aws_iam_role.ec2.name
}
