properties([disableConcurrentBuilds(), buildDiscarder(logRotator(artifactDaysToKeepStr: '5', artifactNumToKeepStr: '5', daysToKeepStr: '5', numToKeepStr: '5'))])

@Library('pipeline-library')
import dk.stiil.pipeline.Constants

podTemplate(yaml: '''
    apiVersion: v1
    kind: Pod
    spec:
      containers:
      - name: kaniko
        image: gcr.io/kaniko-project/executor:v1.20.0-debug
        command:
        - sleep
        args: 
        - 99d
        volumeMounts:
        - name: kaniko-secret
          mountPath: /kaniko/.docker
      - name: golang
        image: golang:1.21.6-alpine3.19
        command:
        - sleep
        args: 
        - 99d
      restartPolicy: Never
      nodeSelector: 
        kubernetes.io/arch: amd64
      volumes:
      - name: kaniko-secret
        secret:
          secretName: dockerhub-dockercred
          items:
          - key: .dockerconfigjson
            path: config.json
''') {
  node(POD_LABEL) {
    String gitCommitMessage
    stage('checkout SCM') {
      scmData = checkout scm
      gitCommitMessage = sh(returnStdout: true, script: "git log --format=%B -n 1 ${scmData.GIT_COMMIT}").trim()
    }
    container('golang') {
      stage('UnitTests') {
        withEnv(['CGO_ENABLED=0', 'GOOS=linux', 'GOARCH=amd64']) {
          sh '''
            go test .
          '''
        }
      }
      stage('Build Application') {
        withEnv(['CGO_ENABLED=0', 'GOOS=linux', 'GOARCH=amd64']) {
          sh '''
            go build -ldflags="-w -s" .
          '''
        }
      }
    }
    if ( !gitCommitMessage.startsWith("renovate/") || ! gitCommitMessage.startsWith("WIP") ) {
      stage('Build Docker Image') {
        container('kaniko') {
          sh '''
            /kaniko/executor --force --context `pwd` --log-format text --destination docker.io/simonstiil/rancher-info-agent:$BRANCH_NAME
          '''
        }
      }
    }
 
  }
}