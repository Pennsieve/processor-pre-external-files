#!groovy

ansiColor('xterm') {
  node('executor') {

  checkout scm

  def authorName  = sh(returnStdout: true, script: 'git --no-pager show --format="%an" --no-patch')
  def serviceName = env.JOB_NAME.tokenize("/")[1]
  def workspace = env.WORKSPACE

  try {
    echo "workspace directory is ${workspace}"

    stage("Run Tests") {
          sh "cd ./service && go test -v ./..."
    }

    stage("Build Container") {
          sh "docker build ."
    }


  } catch (e) {
    slackSend(color: '#b20000', message: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL}) by ${authorName}")
    throw e
  }

  slackSend(color: '#006600', message: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL}) by ${authorName}")
  }
}