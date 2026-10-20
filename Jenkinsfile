pipeline {
    agent any

    stages {
        stage('Format') {
            steps {
                sh 'test -z "$(gofmt -l .)"'
            }
        }
        stage('Vet') {
            steps {
                sh 'go vet ./...'
            }
        }
        stage('Test') {
            steps {
                sh 'go test -v -cover ./...'
            }
        }
        stage('Build') {
            steps {
                sh 'go build -o bin/tracker ./cmd/tracker'
                archiveArtifacts artifacts: 'bin/tracker', fingerprint: true
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished'
        }
    }
}
