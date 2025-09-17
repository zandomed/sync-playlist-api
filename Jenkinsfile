pipeline {
    agent { dockerfile true }

    triggers {
        githubPush()
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
                        image.tag(env.BUILD_NUMBER)
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
        }
        success {
            echo 'Pipeline succeeded! ✅'
        }
        failure {
            echo 'Pipeline failed! ❌'
        }
        unstable {
            echo 'Pipeline is unstable! ⚠️'
        }
    }
}