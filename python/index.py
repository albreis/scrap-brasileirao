import requests
from bs4 import BeautifulSoup
import json

url = "https://www.cbf.com.br/futebol-brasileiro/competicoes/campeonato-brasileiro-serie-a/2023"
response = requests.get(url)
if response.status_code == 200:
    html = response.text
    soup = BeautifulSoup(html, 'html.parser')
    
    # Localize a tabela com a classe "table m-b-20 tabela-expandir"
    table = soup.find('table', class_='table m-b-20 tabela-expandir')
    rows = table.find_all('tr', class_='expand-trigger')

    result = []

    for row in rows:
        # Inicialize um dicionário para cada linha
        teamData = {}

        # Pegue as TDs
        tds = row.find_all('td')
        ths = row.find_all('th')

        if len(tds) > 0:
            # Posição
            positionElement = tds[0].find('b')
            if positionElement:
                teamData['Posição'] = positionElement.text

            # Variação
            variationElement = tds[0].find('span')
            if variationElement:
                teamData['Variação'] = variationElement.text

            # Escudo
            shieldElement = tds[0].find('img')
            teamData['Time'] = {}
            if shieldElement:
                teamData['Time']['Escudo'] = shieldElement['src']

            # Time
            teamElement = tds[0].find_all('span')
            if len(teamElement) > 1:
                teamData['Time']['Nome'] = teamElement[1].text
            
            teamData['PTS'] = ths[0].text
            teamData['Jogos'] = tds[1].text
            teamData['Vitórias'] = tds[2].text
            teamData['Empates'] = tds[3].text
            teamData['Derrotas'] = tds[4].text
            teamData['Gols Pró'] = tds[5].text
            teamData['Gols Contra'] = tds[6].text
            teamData['Saldo de Gols'] = tds[7].text
            teamData['Cartões Amarelo'] = tds[8].text
            teamData['Cartões Vermelho'] = tds[9].text
            teamData['Aproveitamento'] = tds[10].text

            # Recentes
            recentes = [recent.text for recent in tds[11].find_all('span')]
            teamData['Recentes'] = recentes

            # Próximo Jogo
            nextGameElement = tds[12].find('img')
            if nextGameElement:
                teamData['Próximo Jogo'] = {
                    'time': nextGameElement['alt'],
                    'escudo': nextGameElement['src'],
                }

            # Adicione os dados da equipe ao resultado
            result.append(teamData)

    # Exiba o resultado
    print(result)

    # Salve os dados em um arquivo JSON
    with open('dados.json', 'w', encoding='utf-8') as json_file:
        json.dump(result, json_file, ensure_ascii=False, indent=4)
else:
    print(f"Erro ao acessar a página: {response.status_code}")
