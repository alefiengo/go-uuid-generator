pipeline {
    agent any
    
    stages {
        stage('Build') {
            steps {
                script {
                    docker.build("alefiengo/go-uuid-generator:j1", "-f Dockerfile .")
                }
            }
        }
        stage('Push to Registry') {
            steps {
                script {
                    docker.withRegistry("https://registry-1.docker.io/v2/", "credencial-registry-docker") {
                        docker.image("alefiengo/go-uuid-generator:j1").push()
                    }
                }
            }
        }
    }
}