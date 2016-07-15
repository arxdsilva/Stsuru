package main

import (
	"fmt"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
}

// Obs.: O tutorial nao usa DB
// funcao 'save' criada para persistir as informacoes implementadas pela estrutura Page em suas possiveis instancias
// O metodo save recebe p como receptor, que eh um ponteiro para a estrutura Page. Nao recebe parametros e retorna um valor do tipo erro
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// A funcao retorna um erro pois este eh o tipo de retorno de WriteFile (uma biblioteca padrao que escreve pedacos de bites em um arquivo);
// A funcao save retorna o valor do erro e deixa a aplicacao fazer o que precisar com esse erro.
// Caso tudo der certo a funcao retorna 'nil' (zero para ponteiros, interfaces e outros tipos)
// 0600 significa que o arquivo deve ser criado com permissoes ao usuario atual de escrever/ler

// Agora vamos carregar as paginas
// A funcao le o conteudo do arquivo na variavel body e retorna um ponteiro para a 'Pagina', que eh construida com o titulo e corpo (body) fornecidos
func loadPage(title string) *Page {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return &Page{Title: title, Body: body}
}

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a MEGA page.")}
	p1.save()
	p2 := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
