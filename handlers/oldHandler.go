// Obs.: O tutorial nao usa DB
// funcao 'save' criada para persistir as informacoes implementadas pela estrutura Page em suas possiveis instancias
// O metodo save recebe p como receptor, que eh um ponteiro para a estrutura Page. Nao recebe parametros e retorna um valor do tipo erro
// loadPage retorna um erro pois este eh o tipo de retorno de WriteFile (uma biblioteca padrao que escreve pedacos de bites em um arquivo);
// A funcao save retorna o valor do erro e deixa a aplicacao fazer o que precisar com esse erro.
// Caso tudo der certo a funcao retorna 'nil' (zero para ponteiros, interfaces e outros tipos)
// 0600 significa que o arquivo deve ser criado com permissoes ao usuario atual de escrever/ler

// Agora vamos carregar as paginas
// A funcao le o conteudo do arquivo na variavel body e retorna um ponteiro para a 'Pagina', que eh construida com o titulo e corpo (body) fornecidos

package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

const port := ":8080"
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles("index.html"))

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/view/test", http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// saveHandler vai salvar a submissao de formularios localizado nas paginas de edicao.
// Usa a funcao Save para editar o arquivo Test e redirecionar o usuario a pagina inicial
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func mainH() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(port, nil)
	fmt.Println("Server is ONLINE in localhost:", port)
}
