const axios = require('axios');
const { JSDOM } = require('jsdom');
const fs = require('fs');

const url = 'https://www.cbf.com.br/futebol-brasileiro/competicoes/campeonato-brasileiro-serie-a/2023';

async function fetchData() {
  try {
    const response = await axios.get(url);
    const { data } = response;
    
    const { document } = new JSDOM(data).window;
    const rows = document.querySelectorAll('tr.expand-trigger');

    const result = [];

    rows.forEach((row) => {
      const teamData = {};
      const tds = row.querySelectorAll('td');
      const ths = row.querySelectorAll('th');

      if (tds.length > 0) {
        const positionElement = tds[0].querySelector('b');
        if (positionElement) {
          teamData['Posição'] = positionElement.textContent;
        }

        const variationElement = tds[0].querySelector('span');
        if (variationElement) {
          teamData['Variação'] = variationElement.textContent;
        }

        const shieldElement = tds[0].querySelector('img');
        teamData['Time'] = {
          'Escudo': shieldElement ? shieldElement.getAttribute('src') : '',
        };

        const teamElement = tds[0].querySelectorAll('span');
        if (teamElement.length > 1) {
          teamData['Time']['Nome'] = teamElement[1].textContent;
        }

        teamData['PTS'] = ths[0].textContent;
        teamData['Jogos'] = tds[1].textContent;
        teamData['Vitórias'] = tds[2].textContent;
        teamData['Empates'] = tds[3].textContent;
        teamData['Derrotas'] = tds[4].textContent;
        teamData['Gols Pró'] = tds[5].textContent;
        teamData['Gols Contra'] = tds[6].textContent;
        teamData['Saldo de Gols'] = tds[7].textContent;
        teamData['Cartões Amarelo'] = tds[8].textContent;
        teamData['Cartões Vermelho'] = tds[9].textContent;
        teamData['Aproveitamento'] = tds[10].textContent;

        const recentes = [];
        const recentElements = tds[11].querySelectorAll('span');
        recentElements.forEach((recentElement) => {
          recentes.push(recentElement.textContent);
        });
        teamData['Recentes'] = recentes;

        const nextGameElement = tds[12].querySelector('img');
        if (nextGameElement) {
          teamData['Próximo Jogo'] = {
            'time': nextGameElement.getAttribute('alt'),
            'escudo': nextGameElement.getAttribute('src'),
          };
        }

        result.push(teamData);
      }
    });

    // Exiba o resultado
    console.log(result);

    // Salvar os dados em um arquivo JSON
    fs.writeFileSync('dados.json', JSON.stringify(result, null, 2));
  } catch (error) {
    console.error(error);
  }
}

fetchData();
