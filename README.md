# Poligrafo

Validador do contrato do teste [seja-um-guia-back](https://github.com/GuiaBolso/seja-um-guia-back)

# Motivação

Eu valido os testes

E toda vez eu abria na mão e ficava lendo e validando cada uma das regrinhas

Sério, uma ou duas vezes por dia é ate gostoso... mas cansa né (ainda mais pq ngm sabe respeitar contrato neh?!)

Daí eu prometi que faria esse código... e fiz

Você deve estar se perguntando: **"Mas por que você fez o validador do teste Java em GO"**

E a resposta é muito óbiva: **"Pq eu quis"**

Use por sua conta em risco

E não, não adianta reclamar do meu código, eu fiz isso nas horas vagas entre curtir a vida e jogar [Valorant](https://playvalorant.com/pt-br/)... então é feio mesmo

# Como usar

Você tem que ter o `go 1.14+` instalado na sua máquina

Compile e execute o `main`

```sh
/> go build main.go
/> ./main <url do teste rodando>
```

Você pode só rodar, sem compilar

```sh
/> go run main.go -- <url do teste rodando>
```

_(Tô seriamente pensando em fazer um actions pra rodar isso... mas segura a periquita... deixa eu estar animado?!)_

Ah, e como uso Windows (sim, eu uso Windows), você pode até colocar o resultado direto no clipboard, no Powershell (não vamos entrar nessa conversa, ok?!)

```ps
PS /> go run main.go -- <url do teste rodando> | Set-Clipboard
```

# Colaborando

A priori esse código é meu e eu coloquei aqui mais para ter o prazer de falar pra todo mundo que valido testes JVM com Go (e pq eu quis)

Por enquanto, colabore fazendo o teste, candidate-se (não, **não usamos Go**... oficialmente usamos Kotlin)

# Licença

Pode passar
