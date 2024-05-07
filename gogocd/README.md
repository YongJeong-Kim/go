# gocd

## install gocd server
[download gocd](https://www.gocd.org/download)

설치한 경로로 이동하여 server를 실행한다.
```bash
C:\Program Files (x86)\Go Server\bin\start-go-server-service.bat
```
다시 실행(start-go-server-service.bat)했을 때 아래처럼 나오면 실행된 것이다.
![Untitled](https://github.com/YongJeong-Kim/go/assets/30817924/657d62b3-3a0a-4fb1-940c-bdc0db290a76)

## install gocd agent

설치한 경로로 이동하여 agent를 실행한다.
```bash 
C:\Program Files (x86)\Go Agent\bin\start-go-agent-service.bat
```
다시 실행(start-go-agent-service.bat)했을 때 아래처럼 나오면 실행된 것이다.
![Untitled(1)](https://github.com/YongJeong-Kim/go/assets/30817924/73b2efba-13b8-4f6f-889e-cbd34105a9b4)

## 대시보드 접속
- localhost:8153/go
![Untitled(2)](https://github.com/YongJeong-Kim/go/assets/30817924/5051519f-6fbe-44e8-9647-8171437f4ed7)

### Part1: Material
- Material Type
  - Git을 사용할 것이니 Git을 선택
- Repository URL
  - Repository URL 입력
- Advanced Settings
  - Repository Branch(기본 branch가 master로 되어있다.)
    - main(branch 전략에 따라 다를 수 있지만 지금은 main으로 한다.)
  - Username 입력
  - Password 입력

Test Connection을 클릭하여 Connection OK를 확인한다.
![Untitled(3)](https://github.com/YongJeong-Kim/go/assets/30817924/a4311a2f-174c-41f0-adb4-7927cf06468c)

### Part 2: Pipeline Name
![Untitled(4)](https://github.com/YongJeong-Kim/go/assets/30817924/8dffd0c7-2f59-449d-a376-701e174bb84d)
- Pipeline Name
  - 적당히 입력

### Part 3: Stage Details
![Untitled(5)](https://github.com/YongJeong-Kim/go/assets/30817924/d22b2be8-43b3-4e23-a617-16fdf1f658d5)
- Stage Name
  - 적당히 입력

### Part 4: Job and Tasks
![Untitled(6)](https://github.com/YongJeong-Kim/go/assets/30817924/2f2c2452-1beb-4750-942b-576492197bd1)
- Job Name
  - 적당히 입력
- prompt에 실행할 커맨드를 입력한다. 입력 후 엔터를 치면 그 커맨드가 저장된다.
![Untitled(7)](https://github.com/YongJeong-Kim/go/assets/30817924/549380f4-187b-4928-970f-02b45381f7eb)

다음과 같은 실행순서를 갖는다.
```bash 
// 첫번째 실행
$ go test ./...

// 두번째 실행
$ docker compose up -d --build 
```
ci에서 go test ./...를 했음에도 cd에서 또 하는 이유는 ci에서 테스트를 실패해도 push가 된다.
cd에서도 테스트하여 실패하면 배포되지 않도록 하기위함이다.
현재는 main branch에 push했지만 main branch에 바로 push하기 보다 하위 브랜치에서 테스트를 거치고 PR, merge 해야할 것이다. 


Save + Run This  Pipeline
- 파이프라인 저장 후 바로 실행해 볼 수 있다.
Save + Edit Full Config
- 파이프라인 저장 후 상세한 설정화면으로 이동한다.

### DASHBOARD 확인
DASHBOARD를 클릭하면 아래와 같이 pipeline이 생성되었다.
![Untitled(8)](https://github.com/YongJeong-Kim/go/assets/30817924/56cf5d31-5195-4004-83df-d56ad7ce619f)

pipeline의 설정(톱니바퀴)을 클릭하고 TASKS탭으로 이동한다.
![Untitled(9)](https://github.com/YongJeong-Kim/go/assets/30817924/3876e753-1140-49bb-b0aa-487ad61dcea4)

Custom Command를 클릭하여 Working Directory를 수정한다.
![Untitled(10)](https://github.com/YongJeong-Kim/go/assets/30817924/b58d71c0-7c0b-4ea5-bb03-87fe07e7b22c)

현재 repository는 여러 프로젝트가 존재하므로 gogocd로 directory를 이동하여 실행해야 한다. gogocd를 입력하고
Custom Command가 2개 이므로 2개 모두 적용한다.

DASHBOARD로 이동하여 pipeline의 재생버튼을 클릭하여 pipeline을 실행한다.
![Untitled(11)](https://github.com/YongJeong-Kim/go/assets/30817924/51109376-8890-4f13-9962-e4a5f6578dd4)
노란색 바로 변경되어 실행중이며 클릭하여 진행상황을 볼 수 있다.

배포가 완료되면 아래와 같이 초록색 바로 바뀐다.
![Untitled(12)](https://github.com/YongJeong-Kim/go/assets/30817924/df027eae-e05c-4aa2-b78f-165a8e9d77a6)

localhost:8080으로 접속하여 status code 200을 확인해보자

앞으로 main branch에 push할 때마다 pipeline이 실행되며 자동 배포된다.