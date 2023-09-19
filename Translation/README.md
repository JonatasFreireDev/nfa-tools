
# nfa-translation-tool

[Stack Edit](https://stackedit.io/)

## Motivacão

Hoje, o processo para traduzir o aplicativo web para uma nova língua nova é totalmente manual, basicamente recebemos um copy deck e passamos linha a linha por vários arquivos json, alterando as traduções conforme solicitado.

O intuito desde pequeno software é automatizar a **criação**, **leitura do copy deck** e **alteração dos arquivos json**.

## Funcionamento

Ao executar pela primeira vez o programa, dois arquivos são gerados e o programa é fechado. Esses dois arquivos são criados na mesma pasta em que o executável do programa se encontrar.

- **config.json**: Arquivo responsável pela configuração das variáveis do programa.

- **logs.txt**: Arquivo no qual podemos consultar os logs gerados pela aplicação

Assim que o **config.json** for gerado, devemos configura-lo antes de executar novamente o programa.
Esse .json é composto por 5 *keys*:

- **toLocale**: Aqui, informamos o *locale* no qual queremos **criar/editar**. **required**

- **translateFilePath***: Informamos o caminho ate o arquivo **.xlsx**, que contem as *strings* de tradução. **required**

- **copyFromLocale**: Caso o *locale* que vamos alterar, não exista, usamos essa *string* para criar um novo *locale*, pegando todas as strings desse locale ja existente para criar o novo .json.

- **customFilesPaths**: É um array de strings, onde devemos informar os caminhos até os arquivos .json que queremos verificar/alterar isoladamente, fora do contexto de locales.

- **localesPath**: É um array de strings, onde devemos informar os caminhos até as pastas ***"locales"*** de cada repositório.

Com o arquivo de configuração pronto, podemos executar o programa mais uma vez, e se tudo estiver ok, a mágica estará feita !

Caso ocorra algum problema, verifique o arquivo logs.txt, provavelmente o motivo do problema estará registrado lá !

## Features

1. Podemos criar *locales* por cada repositório informado, baseado em um outro *locale* já criado.

2. É possivel consumir uma planilha e alterar varios repos ao mesmo tempo

3. No arquivo de logs, temos a informação de quais arquivos foram modificados, e quais traduções foram substituídas

4. É possivel alterar arquivos .json separadamente, para isso, podemos passar o path do arquivo diretamente no **customFilesPaths**
  

## Observações importantes para o funcionamento

1. Devemos executar o programa a primeira vez para gerar os arquivos, assim que o config for gerado, configurar ele corretamente e executar novamente o programa

2. Infelizmente, a criação de um novo locale, só funciona para o padrão que temos hoje de **locales > en-US > file.json**, qualquer repositório com um padrão diferente desse, não vai ser possível criar ou alterar os arquivos .json utilizando a chave **localesPath**

3. Para alterar arquivos fora do padrao citado a cima, podemos passar o path do arquivo diretamente no **customFilesPaths**

4. A planilha do copy deck, **deve** estar com a **extensão**  **.xlsx**, e deve estar no **seguinte formato**:

	- A1: *String* de referencia (tradução)

	- B1: Nova *string* de tradução

		|  | A   |  B |
		|--|--|--| 
		|1  | add | adicionar  |
		|2  | buy| comprar|


4. Se o arquivo dentro de locale ja estiver criada, a informação **CopyFromLocale** nao será ultilizada.