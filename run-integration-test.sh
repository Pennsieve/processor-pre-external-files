#!/usr/bin/env sh


if [ "$#" -eq 0 ]; then
  echo "must pass an env file. See Readme.md"
  exit 1
fi

for env_file in "$@"; do
    if [ -f "$env_file" ]; then
        set -o allexport && . "$env_file" && set +o allexport
    else
        echo "environment file $env_file is missing"
        exit 1
    fi
done

root_dir=$(pwd)
# the local directory in this directory that docker compose binds to /service/data in the container
local_volume="data"
docker_volume="/service/data"

config_file_name="config.json"
local_config=${root_dir}/${config_file_name}


if [ -z "$INTEGRATION_ID" ]; then
  echo "no value set for INTEGRATION_ID; exiting"
  exit 1
fi

if [ ! -f "$local_config" ]; then
  echo "missing config file $local_config; exiting"
  exit 1
else
  echo "using config file $local_config"
fi

echo "** deleting $root_dir/$local_volume **"
rm -fR "$local_volume"

echo "** creating required directories and config **"
# config
if [[ "$CONFIG_DIR" =~ ^$docker_volume ]]; then
  local_config_dir=${root_dir}/${local_volume}/${CONFIG_DIR#$docker_volume/}
  mkdir -p "$local_config_dir"
  echo "created $local_config_dir"
  cp "$local_config" "$local_config_dir"
  echo "copied $local_config to $local_config_dir"
else
  echo "expected $CONFIG_DIR to start with $docker_volume; exiting"
  exit 1
fi

# input
if [[ "$INPUT_DIR" =~ ^$docker_volume ]]; then
  local_input_dir=${root_dir}/${local_volume}/${INPUT_DIR#$docker_volume/}
  mkdir -p "$local_input_dir"
  echo "created $local_input_dir"
else
  echo "expected $INPUT_DIR to start with $docker_volume; exiting"
  exit 1
fi

# output
if [[ "$OUTPUT_DIR" =~ ^$docker_volume ]]; then
  local_output_dir=${root_dir}/${local_volume}/${OUTPUT_DIR#$docker_volume/}
  mkdir -p "$local_output_dir"
  echo "created $local_output_dir"
else
  echo "expected $OUTPUT_DIR to start with $docker_volume; exiting"
  exit 1
fi

#docker-compose up --build
docker compose run --build --rm pre-processor


