A plugin to generate SBOMs (Software Bill of Materials) using [anchore/syft](https://github.com/anchore/syft).

# Usage

The following settings changes this plugin's behavior.

* select_catalogers (optional) comma-separated list of catalogers to use.
* output (optional) comma-separated list of output formats (e.g. `spdx-json=report.json,cyclonedx-json`).
* source_name (optional) name to use as the source in the SBOM.

The source version is automatically derived from `DRONE_TAG` if set, otherwise `DRONE_COMMIT_SHA`.

Below is an example `.drone.yml` that uses this plugin.

```yaml
kind: pipeline
name: default

steps:
- name: sbom
  image: nemunaire/drone-syft
  pull: if-not-exists
  settings:
    output: spdx-json=sbom.spdx.json
    source_name: my-project
```

Below is an example with multiple outputs and cataloger selection.

```yaml
kind: pipeline
name: default

steps:
- name: sbom
  image: nemunaire/drone-syft
  pull: if-not-exists
  settings:
    select_catalogers: dpkg,rpm
    output: spdx-json=sbom.spdx.json,cyclonedx-json=sbom.cdx.json
    source_name: my-project
```

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t nemunaire/drone-syft -f docker/Dockerfile .
```

# Testing

Execute the plugin from your current working directory:

```text
docker run --rm \
  -e DRONE_COMMIT_SHA=8f51ad7884c5eb69c11d260a31da7a745e6b78e2 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_BUILD_NUMBER=43 \
  -e DRONE_BUILD_STATUS=success \
  -w /drone/src \
  -v $(pwd):/drone/src \
  nemunaire/drone-syft
```
