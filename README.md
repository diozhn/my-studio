# Meu Estúdio 🎨

Uma API simples escrita em Go para artistas postarem seus desenhos.
Idealizada como um presente para minha artista favorita ❤️

---

## ✨ Funcionalidades

- ✅ Upload de imagens com título e legenda
- ✅ Listagem de todas as artes cadastradas
- ✅ Visualização das imagens via rota pública
- ✅ Atualização parcial dos dados (mantendo campos antigos)
- ✅ Remoção de artes com exclusão do arquivo da imagem no disco

---

## 🚀 Tecnologias

- [Go (Golang)](https://golang.org/)
- [Fiber](https://gofiber.io/) – Web framework
- [GORM](https://gorm.io/) – ORM para banco de dados
- SQLite – Banco de dados leve, local

---

## 📦 Como rodar

```bash
#Clonar o projeto
git clone [https://github.com/diozhn/my-studio.git](https://github.com/diozhn/my-studio.git)
cd my-studio

#instalar dependências
go mod tidy

# Rodar o servidor
go run main.go
```
Acesse em:
http://localhost:3000

---

## 📂 Uploads
As imagens são salvas no diretório `uploads/` e podem ser acessadas via URL:

```bash
http://localhost:3000/uploads/nome_da_imagem.jpg
```

---

## 📮 Rotas da API

| Método | Rota            | Descrição                          |
| ------ | --------------- | ---------------------------------- |
| GET    | /artworks     | Lista todas as artes               |
| POST   | /artworks     | Cria uma arte (upload + form-data) |
| PUT    | /artworks/:id | Atualiza campos da arte (JSON)     |
| DELETE | /artworks/:id | Deleta a arte e a imagem do disco  |

---

## 🧪 Exemplo de envio (form-data)

- `title: "Arte linda"`

- `caption: "Feita com carinho"`

- `image: (arquivo de imagem)`

---

## 🚧 Funcionalidades Futuras

- 📸 Postagem simultânea em Instagram, Facebook e Pinterest

- ⭐ Sistema de curtidas/favoritos para as artes

- 🖼️ Galeria web simples para visualização e organização

- 🔐 Autenticação de usuários e controle de permissões

- 📢 Notificações para seguidores sobre novas artes

---

#### Feito com 💙 em Go por @diozhn