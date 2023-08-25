pipeline {
    agent any
        environment {
            registry = "931654140464.dkr.ecr.us-east-1.amazonaws.com/player_service_api"
        }
    options {
        skipStagesAfterUnstable()
    }
    stages {
         stage('Clone repository') {
            steps {
                script{
                checkout scm
                }
            }
        }

        stage('Build') {
            steps {
                script{
                 app = docker.build("player-service-api")
                }
            }
        }
        stage('Deploy') {
            steps {
                script{
                    sh 'aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 931654140464.dkr.ecr.us-east-1.amazonaws.com'
                    sh 'docker push 931654140464.dkr.ecr.us-east-1.amazonaws.com/player_service_api:latest'
                    }
                }
            }
        }
    }
}