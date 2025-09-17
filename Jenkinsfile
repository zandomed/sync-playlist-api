pipeline {
    agent any

    triggers {
        githubPush()
    }

    environment {
        DOCKER_IMAGE_NAME = 'sync-playlist-api'
        DOCKER_REGISTRY = env.DOCKER_REGISTRY ?: 'localhost:5000'
    }

    options {
        githubProjectProperty(projectUrlStr: 'https://github.com/zandomed/sync-playlist-api')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                echo 'Source code checked out successfully'
            }
        }

        stage('Docker Build') {
            steps {
                echo 'Building Docker image...'
                script {
                    def imageTag = "${env.DOCKER_REGISTRY}/${env.DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}"
                    def latestTag = "${env.DOCKER_REGISTRY}/${env.DOCKER_IMAGE_NAME}:latest"

                    def image = docker.build(imageTag, ".")

                    if (env.BRANCH_NAME == 'main') {
                        image.tag('latest')
                        image.tag(latestTag)
                    }

                    env.DOCKER_IMAGE_TAG = imageTag
                }
            }
            post {
                success {
                    echo 'Docker image built successfully'
                }
            }
        }

        // stage('Docker Push') {
        //     when {
        //         anyOf {
        //             branch 'main'
        //             branch 'develop'
        //         }
        //     }
        //     steps {
        //         echo 'Pushing Docker image to registry...'
        //         script {
        //             docker.withRegistry("http://${env.DOCKER_REGISTRY}") {
        //                 def image = docker.image(env.DOCKER_IMAGE_TAG)
        //                 image.push()

        //                 if (env.BRANCH_NAME == 'main') {
        //                     image.push('latest')
        //                 }
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
            script {
                githubNotify(
                    status: 'SUCCESS',
                    description: 'Build completed successfully',
                    context: 'ci/jenkins/build'
                )
            }
            // slackSend(
            //     channel: '#ci-cd',
            //     color: 'good',
            //     message: "✅ Build succeeded for ${env.JOB_NAME} - ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
            // )
        }
        failure {
            echo 'Pipeline failed! ❌'
             script {
                githubNotify(
                    status: 'FAILURE',
                    description: 'Build failed',
                    context: 'ci/jenkins/build'
                )
            }
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