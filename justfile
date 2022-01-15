@_help:
  just --list

# build extension aritfacts
build:
    #!/usr/bin/env bash
    goreleaser build --rm-dist --snapshot
    cd {{justfile_directory()}}/dist && mkdir -p extensions && mv gomplate-lambda-extension_linux_amd64/gomplate-lambda-extension extensions/
    zip -r gomplate-lambda-extension.zip extensions/

region := "us-west-2"

# publish lambda layer
publish: build
    #!/usr/bin/env bash
    arn="$(aws lambda publish-layer-version \
        --region {{region}} \
        --zip-file "fileb://dist/gomplate-lambda-extension.zip" \
        --layer-name "gomplate-lambda-extension" \
        --output text \
        --query LayerVersionArn)"
