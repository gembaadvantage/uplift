# AWS CodeBuild

AWS CodeBuild can be used as a standalone service when running uplift. This guide assumes CodeBuild was configured manually through the AWS Console and only focuses on the gotchas[^1].

CodeBuild will always receive a git clone with a detached HEAD. By default, uplift will [error](../faq/gitdetached.md) in this scenario. When performing a release, this will need to be resolved through a `git checkout`. The `CODEBUILD_SOURCE_VERSION` variable contains the necessary git reference.

## IAM

Additional permissions are needed to pull and push code within AWS CodeBuild. These vary based on the SCM used.

!!!attention "Principle of Least Privilege"

    For illustration purposes, a resource type of `"*"` is used. This should always be narrowed to the specific resource when possible.

### CodeCommit

The `codecommit:GitPush` IAM permission needs to be added. By default, the associated service role will already have the `codecommit:GitPull` permission.

```{ .json .annotate linenums="1" hl_lines="8" }
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "CodeCommitUplift",
      "Effect": "Allow",
      "Action": ["codecommit:GitPush"],
      "Resource": "*"
    }
  ]
}
```

### GitHub

Permissions are granted through the `AWS Connector for GitHub` OAuth application from the AWS Console.

## Buildspec

The buildspec can change depending on the base image used by the CodeBuild project.

### Amazon Images

Tested against the Amazon Linux 2, Ubuntu and Windows variants.

```{ .yaml .annotate linenums="1" hl_lines="5" }
# buildspec.yml

version: 0.2
env:
  git-credential-helper: yes # (1)
phases:
  install:
    commands:
      - curl https://raw.githubusercontent.com/gembaadvantage/uplift/main/scripts/install | bash
  pre_build:
    commands:
      - git checkout ${CODEBUILD_SOURCE_VERSION##"refs/heads/"} # (2)
  build:
    commands:
      - uplift release
```

1. Without this uplift will lack any [credentials](https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec.env.git-credential-helper) when attempting to push code back to the source SCM.
2. This can be simplified to `git checkout $CODEBUILD_SOURCE_VERSION` when cloning from GitHub directly

### Official Uplift Image

Tested against the public `gembaadvantage/uplift` image.

!!!attention "Dealing with DockerHub Rate Limits"

    There are known issues with accessing public DockerHub repositories from AWS services, documented [here](https://aws.amazon.com/blogs/containers/advice-for-customers-dealing-with-docker-hub-rate-limits-and-a-coming-soon-announcement/).

```{ .yaml .annotate linenums="1" hl_lines="5" }
# buildspec.yml

version: 0.2
env:
  git-credential-helper: yes
phases:
  pre_build:
    commands:
      - git checkout ${CODEBUILD_SOURCE_VERSION##"refs/heads/"}
  build:
    commands:
      - uplift release
```

## Clone Depth

While configuring a CodeBuild project, the clone depth can be specified. For simplicity, a full clone should be used. If a shallow clone is preferred, you may need to fetch all tags by using the `--fetch-all` flag.

[^1]: A preferred approach for generating an AWS CodePipeline would be to either write a CloudFormation [template](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-codepipeline-pipeline.html) manually or use the [AWS CDK](https://github.com/aws/aws-cdk) tooling. This is known as Infrastructure as Code (IaC), and wasn't included in the documentation to avoid unnecessary complexity.
