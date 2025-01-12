#!/bin/bash

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
    echo -e "${BLUE}Setting site environment to $1...${NC}"
    terraform output -json \
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
        echo -e "${RED}Unsupported environment${NC}"
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
        *)
            echo -e "${RED}Unsupported subcommand${NC}"
            echo "    $BASH_SOURCE <environment> deploy"
            echo "    $BASH_SOURCE <environment> destroy"
            echo "    $BASH_SOURCE <environment> env"
            echo "    $BASH_SOURCE <environment> output"
            echo
            ;;
    esac
}

format() {
    echo -e "${BLUE}Running formatter...${NC}"
    terraform fmt -recursive
    yamlfmt -dstar **/*.yml
}

output() {
    echo -e "${BLUE}Printing outputs...${NC}"
    terraform output
}

case $1 in
    format)
        format
        ;;
    *)
        environment $1 $2 $3 $4
        ;;
esac
