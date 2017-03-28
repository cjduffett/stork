# Stork

Synthetic patient data as a service, **coming soon!**

## What is Stork?

**Stork** is a cloud-based service that runs instances of **[Synthea](https://github.com/synthetichealth/synthea)**, a synthetic patient record generator. Synthea generates realistic, synthetic patient records that are free from the legal and privacy concerns that plague real healthcare data. This enables software development and experimentation in industry, academia, research, and government that would otherwise not be possible.

## How Do I Use It?

1. Make a request for a large (or small) batch of synthetic patient data.
2. Wait a bit.
3. Check your email for a link to download your newly generated records!

## How Does it Work?

Stork leverages the on-demand scalability of Amazon Web Services, creating instances of Synthea as-needed. Large quantities of Synthea data are generated in-parallel by multiple EC2 instances. All Synthea output is stored in an S3 bucket for download.

## License

Copyright 2017 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

```
http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.