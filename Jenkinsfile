pipeline {
    agent any

    environment {
        REGISTRY = "docker.io/hotpen"
        GO_IMAGE = "${REGISTRY}/sre-go"
        NODE_IMAGE = "${REGISTRY}/sre-node"
        NAMESPACE = "sre-test"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    BRANCH = env.GIT_BRANCH.replace('origin/', '').replace('/', '-')
                    COMMIT = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
                    TAG = "${BRANCH}-${COMMIT}"
                }
            }
        }

        stage('Build Images') {
            steps {
                dir('services/go-service') {
                    sh "docker build -t ${GO_IMAGE}:${TAG} ."
                }
                dir('services/node-service') {
                    sh "docker build -t ${NODE_IMAGE}:${TAG} ."
                }
            }
        }

        stage('Push Images') {
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: 'docker-hub-creds',
                        usernameVariable: 'USER',
                        passwordVariable: 'PASS'
                    )
                ]) {
                    sh 'echo $PASS | docker login -u $USER --password-stdin'
                    sh "docker push ${GO_IMAGE}:${TAG}"
                    sh "docker push ${NODE_IMAGE}:${TAG}"

                    script {
                        if (BRANCH == "main") {
                            sh "docker tag ${GO_IMAGE}:${TAG} ${GO_IMAGE}:latest"
                            sh "docker tag ${NODE_IMAGE}:${TAG} ${NODE_IMAGE}:latest"
                            sh "docker push ${GO_IMAGE}:latest"
                            sh "docker push ${NODE_IMAGE}:latest"
                        }
                    }
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withCredentials([file(credentialsId: 'kubeconfig-file', variable: 'KUBECONFIG_FILE')]) {
                    sh 'mkdir -p $HOME/.kube'
                    sh 'cp $KUBECONFIG_FILE $HOME/.kube/config'
                    sh "kubectl -n sre-test set image deployment/sre-go sre-go=${GO_IMAGE}:${TAG}"
                    sh "kubectl -n sre-test set image deployment/sre-node sre-node=${NODE_IMAGE}:${TAG}"
                }
            }
        }
    }
}
