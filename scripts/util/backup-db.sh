#!/bin/bash

# scripts/util/backup-db.sh
# Author: erumer
# Date: 2025-03-02
# Description: This script uses go-jsonschema to generate a go package for an OSSP artifact schema
# Requirements:
#   - ossp-cli: OSSP CLI
#   - go-jsonschema: JSON Schema go code generator
# Usage: ./json-schema-libs.sh [options]
# Options:
#   -s, --stage: OSSP stage to in which to create the repo
#   -r, --repo: Name of the repo to create
#   -a, --schema: Name of the schema for which to generate lib
#   -x, --version: Version of the schema for which to generate lib
#   -p, --package: Name of the go package for the lib
#   -o, --output: File in which to store the go package go
#   -v, --debug: Print debug information
#   -h, --help: display help

SCRIPT=`basename $0`

usage() {
    printf "Usage: $0 [options]\n"
    printf "\n"
    printf "\t-b, --database \t\t\t\t\tDatabse file.\n"
    printf "\t-o, --output \t\t\t\t\tDatabse file.\n"
    printf "\t-d, --debug \t\t\t\tPrint debug information.\n"
    printf "\t-h, --help \t\t\t\tPrint this usage.\n"
    printf "\n"
}

debug() {
    printf "($SCRIPT) DEBUG: sqlite command = %s\n" "${BACKUP_CMD}"
    printf "\n"
}

# Handle long opts
for arg in "$@"; do
  shift
  case "$arg" in
    '--debug')      set -- "$@" '-d'   ;;
    '--help')       set -- "$@" '-h'   ;;
    *)              set -- "$@" "$arg" ;;
  esac
done


TS=$(date --utc +%FT%TZ)
DB=${db/wartacker.sqlite:-foo}

# Defaults
DEBUG=false
DATABASE="db/wartracker.sqlite3"

# Parse options
OPTIND=1
while getopts ":hds:r:a:x:p:o:" opt; do
    case $opt in
        b)
            if [ -f ${OPTARG} ]; then
                DATABASE=${OPTARG}
            else
                usage
                exit 1
                ;;
            fi
        d)
            DEBUG=true
            ;;
        h)
            usage
            exit 0
            ;;
        \?)
            >&2 printf "($SCRIPT) ERROR: Unsupported option \"-%s\".\n\n" ${OPTARG}
            usage >&2
            exit 1
            ;;
    esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Validate parameters
VALID=true
if [ -z ${REPO} ]; then
    >&2 printf "($SCRIPT) ERROR: Repo not set.\n"
    usage >&2
    VALID=false
fi

if [ -z ${SCHEMA} ]; then
    >&2 printf "($SCRIPT) ERROR: Schema not set.\n"
    usage >&2
    VALID=false
fi

if [ -z ${VERSION} ]; then
    >&2 printf "($SCRIPT) ERROR: Version not set.\n"
    usage >&2
    VALID=false
fi

if [ -z ${PACKAGE} ]; then
    >&2 printf "($SCRIPT) ERROR: Package not set.\n"
    usage >&2
    VALID=false
fi

if [ -z ${OUTPUT} ]; then
    >&2 printf "($SCRIPT) ERROR: Output not set.\n"
    usage >&2
    VALID=false
fi

if [ ${VALID} = false ]; then
    usage >&2
    exit 1
fi

# Configure commands
case ${STAGE} in
    prod)
        OSSP_ENDPOINT="ossp.us-phoenix-1.ocs.oraclecloud.com"
        ;;
    alpha|beta|gamma|dev1|dev2|dev3|dev4|dev5|dev6|dev7|dev8|dev9)
        OSSP_ENDPOINT="${STAGE}.ossp.us-ashburn-1.ocs.oc-test.com"
        ;;
    *)
        >&2 printf "(${SCRIPT}) ERROR: Stage %s is not one of alpha, beta, gamma, or prod\n" ${STAGE}
        usage
        exit 1
        ;;
esac

OSSP_CMD="ossp -z ${OSSP_ENDPOINT} -o json"

GET_SCHEMA_CMD="${OSSP_CMD} get schema ${REPO}:${SCHEMA}:${VERSION}"

if [ ${DEBUG} = true ]; then
    >&2 debug
    >&2 printf "($SCRIPT) DEBUG: No parameter problems detected...\n\n"
fi

# Get OSSP schema
EXISTS=true
GET_SCHEMA_OUTPUT=$(${GET_SCHEMA_CMD})
if [ $? != 0 ]; then
    >&2 printf "($SCRIPT) ERROR: Could not retrieve schema \"%s\".\n" "${REPO}:${SCHEMA}:${VERSION}"
    EXISTS=false
fi

if [ ${DEBUG} = true ]; then
    >&2 printf "\n($SCRIPT) DEBUG: Get Schema OSSP Output --\n%s\n" "${GET_SCHEMA_OUTPUT}"
fi

if [ $EXISTS = false ]; then
    exit 1
fi

# Extract JSON Schema
VALID=true

JSON_SCHEMA=$(echo ${GET_SCHEMA_OUTPUT} | jq '.data.jsonSchema')
if [ $? != 0 ]; then
    >&2 printf "($SCRIPT) ERROR: Could not parse schema \"%s\".\n" "${REPO}:${SCHEMA}:${VERSION}"
    VALID=false
fi

if [ ${DEBUG} = true ]; then
    >&2 printf "\n($SCRIPT) DEBUG: JSON Schema Output --\n%s\n" "${JSON_SCHEMA}"
fi

if [ ${VALID} = false ]; then
    exit 1
fi

# Save JSON Schema
mkdir -p ${SCRATCH_DIR}
echo ${JSON_SCHEMA} > ${SCRATCH_DIR}/${SCHEMA}.json

# Generate lib
GOJSON_OUTPUT=$(go-jsonschema ${SCRATCH_DIR}/${SCHEMA}.json -p ${PACKAGE} -o ${OUTPUT})
if [ $? -ne 0 ]; then
    >&2 printf "($SCRIPT) ERROR: Failed to generate go code for \"%s\"." "${REPO}:${SCHEMA}:${VERSION}"
    if [ ${DEBUG} = true ]; then
        >&2 printf "\n($SCRIPT) DEBUG: go-jsonschema output --\n%s\n" "${GOJSON_OUTPUT}"
    fi
    exit 1
fi

printf "($SCRIPT) INFO: Created go package for schema \"%s\".\n" ${REPO}:${SCHEMA}:${VERSION}

exit 0