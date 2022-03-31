FROM gitpod/workspace-go

USER root

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-374.0.0-linux-x86_64.tar.gz && tar -C /opt/ -xf google-cloud-sdk-374.0.0-linux-x86_64.tar.gz && rm google-cloud-sdk-374.0.0-linux-x86_64.tar.gz

RUN curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH) && chmod +x kubebuilder && mv kubebuilder /bin/

USER gitpod

RUN sudo /opt/google-cloud-sdk/install.sh --usage-reporting false --command-completion true --path-update true --additional-components kubectl --quiet --rc-path /home/gitpod/.bashrc