#!/bin/bash

source config.conf

\cp -f source/kubernetes/$target_system/$kube_version/bin/* ${bin_dir}

