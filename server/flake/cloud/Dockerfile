# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Based on ubuntu
FROM ubuntu:18.04
LABEL maintainers="icebound"
LABEL description="Flake Controller App"

# Copy from build directory
COPY ./cloud-controller /cloud-controller

# Define default command
ENTRYPOINT ["/cloud-controller"]

# Run the executable
CMD ["cloud-controller"]