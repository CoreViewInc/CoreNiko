[![Default](https://github.com/CoreViewInc/CoreNiko/actions/workflows/go.yml/badge.svg)](https://github.com/CoreViewInc/CoreNiko/actions/workflows/go.yml)
![Develop](https://github.com/CoreViewInc/CoreNiko/actions/workflows/go.yml/badge.svg?branch=develop)

# CoreNiko

**CoreNiko** is an innovative project designed to simplify the process of building Docker images within Kubernetes by providing a transparent proxy that harnesses the power of the Kaniko project. CoreNiko allows developers to use familiar Docker commands while taking advantage of Kaniko's advanced features and security model.

### Key Features:

- **Docker CLI Proxy**: CoreNiko serves as a drop-in replacement for the Docker CLI, allowing developers to execute standard Docker build commands that are then transparently translated to utilize Kaniko's executor within a Kubernetes environment.

- **Leveraging Kaniko's Strengths**: Kaniko is trusted for its ability to build container images in Kubernetes without Docker. CoreNiko builds upon this foundation, ensuring a seamless integration with Kubernetes by automating Kaniko's complexities with a Docker-like experience.

- **Enhanced CI/CD Workflows**: By integrating Kaniko's capabilities, CoreNiko enhances CI/CD pipelines with features that improve build efficiency, support advanced caching mechanisms, and ensure that builds are free from potential Docker daemon security concerns.

- **Ease of Transition**: Adopting CoreNiko means there is no need to alter Dockerfiles or learn new workflows. Development teams can switch from Docker to Kaniko-based builds without friction, thanks to CoreNiko's proxy capabilities.

- **Security at the Forefront**: CoreNiko utilizes Kaniko's secure methodology for building Docker images to avoid the security pitfalls associated with running a Docker daemon, particularly in cloud-native environments.

### How it Works:

CoreNiko acts as an intermediary proxy that simplifies the use of Kaniko. Developers can perform Docker builds using the usual commands:

```sh
$ docker build -t my-image .
```

In the background, CoreNiko processes this command and leverages Kaniko to perform the actual image build within a Kubernetes cluster. The proxy ensures that the Docker command's functionality is preserved without the need for direct interaction with Kaniko's command structure.

### Contributing to CoreNiko:

We welcome contributions to the CoreNiko project! Whether you're fixing bugs, adding features, improving documentation, or helping the community, your efforts will make a significant impact. Here's how you can contribute:

#### Reporting Issues:

1. Utilize the project's issue tracker to report bugs or suggest enhancements.
2. Provide detailed descriptions and reproducible steps if you're reporting a bug.

#### Code Contributions:

1. Fork the repository and create a new branch for your contribution.
2. Follow the code style and contribution guidelines of the project.
3. Update or create tests as necessary.
4. Ensure your contributions pass existing tests.
5. Commit your changes with clear and understandable messages.
6. Push your changes and open a pull request against the original repo.

#### Community Engagement:

- Provide assistance in community support channels.
- Share your CoreNiko use cases to guide and inspire others.

With your help, we can continue building upon the success of the Kaniko project to make CoreNiko a robust and user-friendly solution for Kubernetes-native Docker image builds.

### Conclusion:

CoreNiko is a strategic tool that aligns with contemporary DevOps practices by combining the Docker CLI's approachability with Kaniko's security and effectiveness in Kubernetes environments. Through CoreNiko, developers gain the best of both worlds, with a secure, efficient, and Kubernetes-optimized image building process that feels familiar and integrates effortlessly into their existing workflows.