require 'net/http'
require 'nokogiri'
require 'json'

url = 'https://www.cbf.com.br/futebol-brasileiro/competicoes/campeonato-brasileiro-serie-a/2023'
uri = URI(url)
response = Net::HTTP.get(uri)

doc = Nokogiri::HTML(response)
result = []

# Localize a tabela com a classe "table m-b-20 tabela-expandir"
table = doc.at('table.table.m-b-20.tabela-expandir')
rows = table.css('tr.expand-trigger')

rows.each do |row|
  team_data = {}

  # Pegue as colunas da linha
  columns = row.css('td')
  headers = row.css('th')

  if columns.any?
    # Posição e Variação
    position_element = columns[0].at('b')
    variation_element = columns[0].at('span')

    team_data['Posição'] = position_element.text.strip if position_element
    team_data['Variação'] = variation_element.text.strip if variation_element

    # Escudo e Nome do Time
    shield_element = columns[0].at('img')
    team_element = columns[0].css('span')[1]

    if shield_element
      team_data['Time'] = {
        'Escudo' => shield_element['src'],
        'Nome' => team_element.text.strip
      }
    end

    # Outros dados da equipe
    team_data['PTS'] = headers[0].text.strip
    team_data['Jogos'] = columns[1].text.strip
    team_data['Vitórias'] = columns[2].text.strip
    team_data['Empates'] = columns[3].text.strip
    team_data['Derrotas'] = columns[4].text.strip
    team_data['Gols Pró'] = columns[5].text.strip
    team_data['Gols Contra'] = columns[6].text.strip
    team_data['Saldo de Gols'] = columns[7].text.strip
    team_data['Cartões Amarelo'] = columns[8].text.strip
    team_data['Cartões Vermelho'] = columns[9].text.strip
    team_data['Aproveitamento'] = columns[10].text.strip

    # Recentes
    recentes = columns[11].css('span').map { |e| e.text.strip }
    team_data['Recentes'] = recentes

    # Próximo Jogo
    next_game_element = columns[12].at('img')

    if next_game_element
      team_data['Próximo Jogo'] = {
        'time' => next_game_element['alt'],
        'escudo' => next_game_element['src']
      }
    end

    result << team_data
  end
end

# Exiba o resultado
puts JSON.pretty_generate(result)

# Salve os dados em um arquivo JSON
File.write('dados.json', JSON.pretty_generate(result))
