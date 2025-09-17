pipeline {
    agent any

    tools {
        go 'go-1.25'
    }
    
    triggers {
        githubPush()
    }

    environment {
        GO111MODULE = 'on'
        GOPROXY = 'https://proxy.golang.org,direct'
        GOSUMDB = 'sum.golang.org'
        CGO_ENABLED = '0'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                echo 'Source code checked out successfully'
            }
        }

        stage('Dependencies') {
            steps {
                echo 'Installing dependencies...'
                sh 'make deps'
            }
        }

        stage('Format Check') {
            steps {
                echo 'Checking code formatting...'
                script {
                    def formatResult = sh(script: 'make format', returnStatus: true)
                    if (formatResult != 0) {
                        error('Code formatting check failed. Please run "make format" locally and commit the changes.')
                    }
                }
            }
        }

        stage('Lint') {
            steps {
                echo 'Running linter...'
                sh 'make lint'
            }
            post {
                always {
                    // publishHTML([
                    //     allowMissing: false,
                    //     alwaysLinkToLastBuild: true,
                    //     keepAll: true,
                    //     reportDir: '.',
                    //     reportFiles: 'golangci-lint-report.xml',
                    //     reportName: 'Lint Report'
                    // ])
                }
            }
        }

        // stage('Test') {
        //     steps {
        //         echo 'Running tests...'
        //         sh 'make test-coverage'
        //     }
        //     post {
        //         always {
        //             // publishHTML([
        //             //     allowMissing: false,
        //             //     alwaysLinkToLastBuild: true,
        //             //     keepAll: true,
        //             //     reportDir: '.',
        //             //     reportFiles: 'coverage.html',
        //             //     reportName: 'Coverage Report'
        //             // ])
        //         }
        //     }
        // }

        stage('Build') {
            steps {
                echo 'Building application...'
                sh 'make build'
            }
            post {
                success {
                    archiveArtifacts artifacts: 'dist/*', fingerprint: true
                }
            }
        }

        // stage('Docker Build') {
        //     when {
        //         anyOf {
        //             branch 'main'
        //             branch 'develop'
        //             changeRequest()
        //         }
        //     }
        //     steps {
        //         echo 'Building Docker image...'
        //         script {
        //             def image = docker.build("sync-playlist-api:${env.BUILD_NUMBER}")
        //             if (env.BRANCH_NAME == 'main') {
        //                 image.tag('latest')
        //             }
        //         }
        //     }
        // }
    }

    post {
        always {
            echo 'Pipeline execution completed'
            cleanWs()
        }
        success {
            echo 'Pipeline succeeded! ✅'
            // slackSend(
            //     channel: '#ci-cd',
            //     color: 'good',
            //     message: "✅ Build succeeded for ${env.JOB_NAME} - ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
            // )
        }
        failure {
            echo 'Pipeline failed! ❌'
            // slackSend(
            //     channel: '#ci-cd',
            //     color: 'danger',
            //     message: "❌ Build failed for ${env.JOB_NAME} - ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
            // )
        }
        unstable {
            echo 'Pipeline is unstable! ⚠️'
            // slackSend(
            //     channel: '#ci-cd',
            //     color: 'warning',
            //     message: "⚠️ Build unstable for ${env.JOB_NAME} - ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
            // )
        }
    }
}