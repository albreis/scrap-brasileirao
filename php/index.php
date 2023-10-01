<?php

$url = "https://www.cbf.com.br/futebol-brasileiro/competicoes/campeonato-brasileiro-serie-a/2023";
$ch = curl_init($url);
curl_setopt_array($ch, [
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_SSL_VERIFYPEER => false
]);
$response = curl_exec($ch);
$error = curl_error($ch);
curl_close($ch);
if ($error) {
    print_r($error);
} else {
    // Desabilitar erros do DomDocument
    libxml_use_internal_errors(true);
    $dom = new DOMDocument();
    $dom->loadHTML($response);
    // Localize a tabela com a classe "table m-b-20 tabela-expandir"
    $xpath = new DOMXPath($dom);
    $rows = $xpath->query('//tr[@class="expand-trigger"]');

    foreach ($rows as $row) {
        // Inicialize um array para cada linha
        $teamData = array();

        // Pegue a primeira TD
        $tds = $row->getElementsByTagName('td');
        $ths = $row->getElementsByTagName('th');

        if ($tds->length > 0) {
            // Posição
            $positionElement = $tds->item(0)->getElementsByTagName('b');
            if ($positionElement->length > 0) {
                $teamData['Posição'] = $positionElement->item(0)->textContent;
            }

            // Variação
            $variationElement = $tds->item(0)->getElementsByTagName('span');
            if ($variationElement->length > 0) {
                $teamData['Variação'] = $variationElement->item(0)->textContent;
            }

            // Escudo
            $shieldElement = $tds->item(0)->getElementsByTagName('img');
            $teamData['Time'] = [];
            if ($shieldElement->length > 0) {
                $teamData['Time']['Escudo'] = $shieldElement->item(0)->getAttribute('src');
            }

            // Time
            $teamElement = $tds->item(0)->getElementsByTagName('span');
            if ($teamElement->length > 0) {
                $teamData['Time']['Nome'] = $teamElement->item(1)->textContent;
            }
            
            $teamData['PTS'] = $ths->item(0)->textContent;
            $teamData['Jogos'] = $tds->item(1)->textContent;
            $teamData['Vitórias'] = $tds->item(2)->textContent;
            $teamData['Empates'] = $tds->item(3)->textContent;
            $teamData['Derrotas'] = $tds->item(4)->textContent;
            $teamData['Gols Pró'] = $tds->item(5)->textContent;
            $teamData['Gols Contra'] = $tds->item(6)->textContent;
            $teamData['Saldo de Gols'] = $tds->item(7)->textContent;
            $teamData['Cartões Amarelo'] = $tds->item(8)->textContent;
            $teamData['Cartões Vermelho'] = $tds->item(9)->textContent;
            $teamData['Aproveitamento'] = $tds->item(10)->textContent;


            // Recentes
            $recentes = array();
            $recentElements = $tds->item(11)->getElementsByTagName('span');
            foreach ($recentElements as $recentElement) {
                $recentes[] = $recentElement->textContent;
            }
            $teamData['Recentes'] = $recentes;

            // Próximo Jogo
            $nextGameElement = $tds->item(12)->getElementsByTagName('img');
            if ($nextGameElement->length > 0) {
                $teamData['Próximo Jogo'] = [
                    'time' => $nextGameElement->item(0)->getAttribute('alt'),
                    'escudo' => $nextGameElement->item(0)->getAttribute('src'),
                ];
            }
        }

        // Adicione os dados da equipe ao resultado
        $result[] = $teamData;
    }
    // Exiba o resultado
    print_r($result);

    file_put_contents('dados.json', json_encode($result, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES));
}
