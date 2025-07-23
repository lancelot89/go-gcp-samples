FROM openjdk:21-slim
ENV GCLOUD_SDK_VERSION=460.0.0
RUN apt-get update && apt-get install -y curl python3 &&     curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_SDK_VERSION}-linux-x86_64.tar.gz &&     tar -xf google-cloud-sdk-${GCLOUD_SDK_VERSION}-linux-x86_64.tar.gz &&     ./google-cloud-sdk/install.sh --quiet &&     rm google-cloud-sdk-${GCLOUD_SDK_VERSION}-linux-x86_64.tar.gz
ENV PATH /google-cloud-sdk/bin:$PATH
RUN gcloud components install beta cloud-firestore-emulator --quiet
CMD ["gcloud", "beta", "emulators", "firestore", "start", "--host-port=0.0.0.0:8080", "--quiet"]
