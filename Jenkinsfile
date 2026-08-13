pipeline {
    agent any

    environment {
        GIT_REPO_URL = 'https://github.com/noctislucius20/VeggieMart.git'
        GIT_BRANCH = 'main'
        GIT_CREDENTIALS_ID = 'github-credential'
        COMPOSE_FILE = 'docker-compose.yml'
        DOCKER_BUILDKIT = '1'
    }

    stages {
        stage('Clone Repository') {
            steps {
                script {
                    try {
                        git(
                            url: "${env.GIT_REPO_URL}",
                            branch: "${env.GIT_BRANCH}",
                            credentialsId: "${env.GIT_CREDENTIALS_ID}"
                        )
                        echo 'Repository cloned successfully.'
                    } catch (err) {
                    error "Failed to clone repository: ${err}"
                    }

                }
            }
        }

        stage('Docker Compose Down') {
            steps {
                script {
                    try {
                        sh """ 
                            cd api
                            docker compose -f ${env.COMPOSE_FILE} down || true
                        """
                        echo 'Old containers stopped successfully.'
                    } catch (err) {
                        echo "Warning: Failed to stop old containers: ${err}"
                    }
                }
            }
        }

        stage('Generate .env file') {
            steps {
               withCredentials([
                    file(credentialsId: 'veggiemart-user-env', variable: 'VMART_USER_ENV'),
                    file(credentialsId: 'veggiemart-product-env', variable: 'VMART_PRODUCT_ENV'),
                    file(credentialsId: 'veggiemart-order-env', variable: 'VMART_ORDER_ENV'),
                    file(credentialsId: 'veggiemart-payment-env', variable: 'VMART_PAYMENT_ENV'),
                    file(credentialsId: 'veggiemart-notification-env', variable: 'VMART_NOTIF_ENV'),
                    file(credentialsId: 'veggiemart-api-gateway-env', variable: 'VMART_API_GATEWAY_ENV'),
                    file(credentialsId: 'veggiemart-docker-env', variable: 'VMART_DOCKER_ENV')
                ]) {

                    sh '''
                        cp "$VMART_USER_ENV" api/user-service/.env
                        cp "$VMART_PRODUCT_ENV" api/product-service/.env
                        cp "$VMART_ORDER_ENV" api/order-service/.env
                        cp "$VMART_PAYMENT_ENV" api/payment-service/.env
                        cp "$VMART_NOTIF_ENV" api/notification-service/.env
                        cp "$VMART_API_GATEWAY_ENV" api/api-gateway/.env
                        cp "$VMART_DOCKER_ENV" api/.env
                    '''
                }
            }
        }

        stage('Docker Compose Build') {
            steps {
                script {
                    try {
                        sh """
                            cd api
                            docker compose -f ${env.COMPOSE_FILE} build
                        """
                        echo 'Docker images built successfully.'
                    } catch (err) {
                        error "Failed to build docker images: ${err}"
                    }
                }
            }
        }

        stage('Docker Compose Up') {
            steps {
                script {
                    try {
                        sh """
                            cd api
                            docker compose -f ${env.COMPOSE_FILE} up -d
                        """
                        echo 'Containers started successfully.'
                    } catch (err) {
                        error "Failed to start containers: ${err}"
                    }
                }
            }
        }
    }

    post {
        success {
            echo '✅ Deployment Success!'
        }

        failure {
            echo '❌ Deployment Failed!'
        }

        cleanup {
            echo '�� Cleaning workspace...'
            cleanWs()
        }
    }
}