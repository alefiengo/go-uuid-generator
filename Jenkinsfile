pipeline {
    agent any
    
    stages {
        stage('Build') {
            steps {
                script {
                    docker.build("go-uuid-generator:j1", "-f Dockerfile .")
                }
            }
        }
    }
}