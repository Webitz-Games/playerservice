pipeline {

    agent any

    environment {
        registry = "931654140464.dkr.ecr.us-east-1.amazonaws.com/player_service_api"
    }
    stages {

        stage ('Checkout') {
            steps {
                checkout scm
            }
        }
        stage ('Docker Build') {
            steps {
                script {
                    dockerImage = docker.build("${env.registry}")
                }
            }
        }
        stage ("Docker Push") {
            steps {
                script {
                    sh 'aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 931654140464.dkr.ecr.us-east-1.amazonaws.com'
                    sh 'docker push 931654140464.dkr.ecr.us-east-1.amazonaws.com/player_service_api:latest'
                }
            }
        }
    }
}