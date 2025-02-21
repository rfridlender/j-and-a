#!/bin/bash

API_DIRECTORY_NAME="api"
ENVIRONMENTS_DIRECTORY_NAME="environments"
SITE_DIRECTORY_NAME="site"
WORKING_DIRECTORY_NAME=$(pwd)

RED="\033[1;31m"
GREEN="\033[1;32m"
YELLOW="\033[1;33m"
BLUE="\033[1;34m"
NC="\033[0m"

if [ -e .env ];
then
    set -o allexport

    source .env set
fi

if [ -e .env.tokens ];
then
    set -o allexport

    source .env.tokens set
fi

build() {
    for dir in cmd/*; do
    if [ -d $dir ]; then
        echo "Building $(basename $dir)..."
        GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o $dir/bootstrap $dir/main.go
        echo "Build complete for $(basename $dir)."
    fi
    done
}

deploy() {
    terraform init -input=false
    echo

    echo -e "${BLUE}Applying $1 infrastructure...${NC}"
    terraform apply -input=false -auto-approve
    echo
}

destroy() {
    echo -e "${BLUE}Destroying $1 environment...${NC}"
    terraform destroy
}

env() {
    terraform_outputs=$(terraform output -json)

    echo -e "${BLUE}Setting project environment to $1...${NC}"
    echo $terraform_outputs \
        | jq -r 'to_entries[] | "\(.key | ascii_upcase)=\(.value.value)"' \
        > "../../.env"

    echo -e "${BLUE}Setting site environment to $1...${NC}"
    echo $terraform_outputs \
        | jq -r 'to_entries[] | select(.key | startswith("vite_")) | "\(.key | ascii_upcase)=\(.value.value)"' \
        > "../../$SITE_DIRECTORY_NAME/.env.local"
}

environment() {
    is_valid_environment=false
    for dir in $ENVIRONMENTS_DIRECTORY_NAME/*;
    do
        if [[ $dir == "$ENVIRONMENTS_DIRECTORY_NAME/$1" ]];
        then
            is_valid_environment=true
            break
        fi
    done

    if [[ $is_valid_environment == false ]]
    then
        echo -e "${RED}Unsupported environment.${NC}"
        echo -e "${BLUE}Please choose from the following environments:${NC}"
        i=1
        environment_directories=()
        for dir in $ENVIRONMENTS_DIRECTORY_NAME/*;
        do
            echo "    $i) ${dir##"$ENVIRONMENTS_DIRECTORY_NAME/"}"
            ((i++))
            environment_directories+=($dir)
        done

        read new_environment_index
        ((new_environment_index--))
        if [[ $new_environment_index -ge ${#environment_directories[@]} ]];
        then
            echo -e "${BLUE}Exiting...${NC}"
            exit 1
        fi
        source $0 ${environment_directories[$new_environment_index]##"$ENVIRONMENTS_DIRECTORY_NAME/"} $2 $3 $4
        exit 1
    fi

    echo -e "${BLUE}Switching to $1 environment...${NC}"
    cd "$ENVIRONMENTS_DIRECTORY_NAME/$1"

    case $2 in
        deploy)
            deploy $1
            ;;
        destroy)
            destroy $1
            ;;
        env)
            env $1
            ;;
        output)
            output
            ;;
        user-create)
            user-create $3
            ;;
        user-delete)
            user-delete $3
            ;;
        user-refresh)
            user-refresh $3
            ;;
        *)
            echo -e "${RED}Unsupported subcommand.${NC}"
            echo "    $BASH_SOURCE <environment> deploy"
            echo "    $BASH_SOURCE <environment> destroy"
            echo "    $BASH_SOURCE <environment> env"
            echo "    $BASH_SOURCE <environment> output"
            echo "    $BASH_SOURCE <environment> user-create <email>"
            echo "    $BASH_SOURCE <environment> user-delete <email>"
            echo "    $BASH_SOURCE <environment> user-refresh"
            echo
            ;;
    esac
}

format() {
    echo -e "${BLUE}Running formatter...${NC}"
    terraform fmt -recursive
    yamlfmt -dstar **/*.yml
    go fmt ./...
    cd $SITE_DIRECTORY_NAME && npm run format
}

output() {
    echo -e "${BLUE}Printing outputs...${NC}"
    terraform output
}

user-create() {
    echo -e "${BLUE}Creating user $1...${NC}"
    aws cognito-idp admin-create-user \
        --user-pool-id $USER_POOL_ID \
        --username $1 \
        --user-attributes \
            Name=email,Value=$1 \
            Name=email_verified,Value=True \
            Name=given_name,Value=Dummy \
            Name=family_name,Value=User \
        --no-cli-pager

    echo -e "${BLUE}Initiating auth...${NC}"
    session=$(
        aws cognito-idp admin-initiate-auth \
            --user-pool-id $USER_POOL_ID \
            --client-id $USER_POOL_CLIENT_ID \
            --auth-flow USER_AUTH \
            --auth-parameters USERNAME=$1,PREFERRED_CHALLENGE=EMAIL_OTP \
            --no-cli-pager \
            | jq -r '.Session'
    )

    echo -e "${BLUE}Please enter the verification code sent to $1:${NC}"
    read verification_code

    echo
    echo -e "${BLUE}Generating tokens...${NC}"
    tokens=$(
        aws cognito-idp admin-respond-to-auth-challenge \
            --user-pool-id $USER_POOL_ID \
            --client-id $USER_POOL_CLIENT_ID \
            --challenge-name EMAIL_OTP \
            --challenge-responses USERNAME=$1,EMAIL_OTP_CODE=$verification_code \
            --session $session \
            --no-cli-pager
    )

    echo -e "${BLUE}Writing tokens to .env.tokens...${NC}"
    echo $tokens \
        | jq -r '.AuthenticationResult.RefreshToken | "REFRESH_TOKEN=\(.)"' \
        > "../../.env.tokens"

    echo -e "${BLUE}Writing tokens to $API_DIRECTORY_NAME/.env...${NC}"
    echo $tokens \
        | jq -r '.AuthenticationResult | to_entries[] | select(.key | endswith("Token")) | "\(.key | gsub("(?<x>[a-z])Token"; "\(.x)_Token") | ascii_upcase)=\(.value)"' \
        > "../../$API_DIRECTORY_NAME/.env"
    echo "API_ENDPOINT=$API_ENDPOINT" >> "../../$API_DIRECTORY_NAME/.env"
}

user-delete() {
    echo -e "${BLUE}Deleting user $1...${NC}"
    aws cognito-idp admin-delete-user \
        --user-pool-id $USER_POOL_ID \
        --username $1
}

user-refresh() {
    echo -e "${BLUE}Refreshing tokens...${NC}"

    tokens=$(
        aws cognito-idp initiate-auth \
            --client-id $USER_POOL_CLIENT_ID \
            --auth-flow REFRESH_TOKEN_AUTH \
            --auth-parameters REFRESH_TOKEN=$REFRESH_TOKEN
    )

    echo -e "${BLUE}Writing tokens to $API_DIRECTORY_NAME/.env...${NC}"
    echo $tokens \
        | jq -r '.AuthenticationResult | to_entries[] | select(.key | endswith("Token")) | "\(.key | gsub("(?<x>[a-z])Token"; "\(.x)_Token") | ascii_upcase)=\(.value)"' \
        > "../../$API_DIRECTORY_NAME/.env"
    echo "API_ENDPOINT=$API_ENDPOINT" >> "../../$API_DIRECTORY_NAME/.env"
}

case $1 in
    build)
        build
        ;;
    format)
        format
        ;;
    *)
        environment $1 $2 $3 $4
        ;;
esac
