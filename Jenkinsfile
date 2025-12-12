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
                    env.BRANCH = env.GIT_BRANCH.replace('origin/', '').replace('/', '-')
                    env.COMMIT = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
                    env.TAG = "${env.BRANCH}-${env.COMMIT}"
                }
            }
        }

        stage('Go Tests') {
            steps {
                dir('services/go-service') {
                    sh 'go test ./... -coverprofile=coverage.out'
                }
            }
        }

        stage('Node Tests') {
            steps {
                dir('services/node-service') {
                    sh 'npm install'
                    sh 'npm test -- --coverage'
                }
            }
        }

        stage('SonarQube Analysis') {
            steps {
                withSonarQubeEnv('sonarqube-server') {
                    sh """
                        sonar-scanner \
                        -Dsonar.projectKey=sre-test \
                        -Dsonar.sources=. \
                        -Dsonar.go.coverage.reportPaths=services/go-service/coverage.out \
                        -Dsonar.javascript.lcov.reportPaths=services/node-service/coverage/lcov.info
                    """
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
            when {
                branch "main"
            }
            steps {
                withCredentials([
                    file(credentialsId: 'kubeconfig-file', variable: 'KUBECONFIG_FILE')
                ]) {

                    sh '''
                        mkdir -p $HOME/.kube
                        cp $KUBECONFIG_FILE $HOME/.kube/config
                    '''

                    script {
                        env.IMAGE_TAG = TAG

                        sh """
                            echo "Deploying Go service with tag ${IMAGE_TAG}"
                            export IMAGE_TAG=${IMAGE_TAG}
                            envsubst < kubernetes-manifest/go-deployment.yaml | kubectl apply -f -

                            echo "Deploying Node service with tag ${IMAGE_TAG}"
                            envsubst < kubernetes-manifest/node-deployment.yaml | kubectl apply -f -
                        """
                    }
                }
            }
        }
    }
}
