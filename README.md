# Meu Estúdio 🎨

Uma API escrita em Go para artistas postarem, visualizarem e interagirem com artes digitais.  
Idealizada como um presente para minha artista favorita ❤️

---

## ✨ Funcionalidades

- Upload de imagens com título e legenda
- Listagem de todas as artes cadastradas, com filtros e paginação
- Visualização das imagens via rota pública e galeria web simples
- Atualização e remoção de artes (com exclusão do arquivo da imagem)
- Sistema de curtidas para as artes
- Filtros por autor, data e título
- Autenticação de usuários (JWT)
- Atualização de perfil de usuário
- Listagem de artes de um usuário específico

---

## 🚀 Tecnologias

- [Go (Golang)](https://golang.org/)
- [Fiber](https://gofiber.io/) – Web framework
- [GORM](https://gorm.io/) – ORM para banco de dados
- [Supabase](https://supabase.com/) – Banco de dados PostgreSQL gerenciado (pode ser substituído por outro PostgreSQL)
- [godotenv](https://github.com/joho/godotenv) – Carregamento de variáveis de ambiente

---

## 📦 Como rodar

```bash
# Clonar o projeto
git clone https://github.com/diozhn/my-studio.git
cd my-studio

# Instalar dependências
go mod tidy

# Copiar o exemplo de variáveis de ambiente e editar
cp .env.example .env
# Edite o .env com suas credenciais do Supabase/PostgreSQL e JWT_SECRET

# Rodar o servidor
go run main.go
```
Acesse em:
http://localhost:3000

---

## 📂 Uploads
As imagens são salvas no diretório `uploads/` e podem ser acessadas via URL:

```
http://localhost:3000/uploads/nome_da_imagem.jpg
```

---

## 📮 Rotas da API

### Autenticação
| Método | Rota             | Descrição                       |
| ------ | ---------------- | ------------------------------- |
| POST   | /register        | Cria um novo usuário            |
| POST   | /login           | Realiza login e retorna tokens  |
| POST   | /refresh-token   | Gera novo token de acesso       |

### Usuários
| Método | Rota                        | Descrição                                 |
| ------ | --------------------------- | ----------------------------------------- |
| GET    | /users/:id                  | Busca perfil de usuário                   |
| PATCH  | /users/:id                  | Atualiza perfil do usuário (autenticado)  |
| GET    | /users/:id/artworks         | Lista artes de um usuário                 |

### Artes
| Método | Rota                        | Descrição                                 |
| ------ | --------------------------- | ----------------------------------------- |
| GET    | /artworks                   | Lista todas as artes (com filtros)        |
| POST   | /artworks                   | Cria uma arte (upload + form-data, auth) |
| GET    | /artworks/:id               | Busca arte por ID                         |
| PATCH  | /artworks/:id               | Atualiza arte (autenticado e dono)        |
| DELETE | /artworks/:id               | Deleta a arte (autenticado e dono)        |
| POST   | /artworks/:id/like          | Curte uma arte                            |
| GET    | /top-artworks               | Lista artes mais curtidas                 |
| GET    | /gallery                    | Galeria web simples (HTML)                |
| GET    | /artworks/filter            | Lista artes filtradas                     |

---

## 🛡️ Segurança

- JWT_SECRET e DATABASE_URL devem ser definidos em variáveis de ambiente (.env)
- Nunca exponha segredos ou senhas publicamente
- Recomenda-se uso de HTTPS em produção

---

## 🗄️ Banco de Dados

- Utiliza PostgreSQL (Supabase recomendado para facilidade e plano gratuito)
- Migrations automáticas via GORM ao iniciar a aplicação

---

## 🧪 Exemplo de envio (form-data)

- `title: "Arte linda"`
- `caption: "Feita com carinho"`
- `image: (arquivo de imagem)`

---

#### Feito com 💙 em Go por @diozhn