ARG UBUNTU_VERSION=24.04
FROM ubuntu:${UBUNTU_VERSION}

ARG PUBLIC_KEY

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        openssh-server sudo ca-certificates curl bzip2 \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /run/sshd \
    && useradd --create-home --shell /bin/bash deploy \
    && printf 'deploy ALL=(ALL) NOPASSWD:ALL\n' > /etc/sudoers.d/deploy \
    && chmod 0440 /etc/sudoers.d/deploy

RUN mkdir -p /home/deploy/.ssh \
    && printf '%s\n' "$PUBLIC_KEY" > /home/deploy/.ssh/authorized_keys \
    && chmod 700 /home/deploy/.ssh \
    && chmod 600 /home/deploy/.ssh/authorized_keys \
    && chown -R deploy:deploy /home/deploy/.ssh

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D", "-e"]
