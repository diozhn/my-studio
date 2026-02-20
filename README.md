## My Studio (TypeScript) 🎨

Uma API agora escrita em **Node.js + TypeScript** para artistas postarem, visualizarem e interagirem com artes digitais.  
Originalmente o projeto foi feito em Go, mas toda a API foi refatorada para TypeScript.

---

## ✨ Funcionalidades

- Upload de imagens com título e legenda
- Listagem de todas as artes cadastradas, com filtros e paginação
- Visualização das imagens via rota pública e galeria web simples
- Atualização e remoção de artes (com exclusão do arquivo da imagem)
- Sistema de curtidas para as artes
- Filtros por autor, data e título
- Autenticação de usuários com JWT (login, registro e refresh token)
- Atualização de perfil de usuário
- Listagem de artes de um usuário específico

> Observação: o fluxo de OAuth social (Google, Instagram, Twitter) do código Go **ainda não foi reimplementado** em TypeScript.

---

## 🚀 Tecnologias

- Node.js + TypeScript
- Express – Web framework
- PostgreSQL – via driver `pg`
- Multer – upload de arquivos
- jsonwebtoken – autenticação JWT
- bcrypt – hash de senha
- dotenv – variáveis de ambiente

---

## 📦 Como rodar (TypeScript)

```bash
# Clonar o projeto
git clone https://github.com/diozhn/my-studio.git
cd my-studio

# Instalar dependências
npm install

# Criar arquivo de variáveis de ambiente
cp .env.example .env
# Edite o .env com suas credenciais do PostgreSQL e JWT_SECRET

# Compilar o TypeScript
npm run build

# Rodar o servidor compilado
npm start

# ou em modo desenvolvimento (ts-node-dev)
npm run dev
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

## 📮 Rotas da API (versão TypeScript)

### Autenticação
| Método | Rota           | Descrição                              |
| ------ | -------------- | -------------------------------------- |
| POST   | /register      | Cria um novo usuário                   |
| POST   | /login         | Realiza login e retorna tokens         |
| POST   | /refresh-token | Gera novo token de acesso              |

**Observações:**
- Senhas são armazenadas com bcrypt.
- O refresh token é armazenado no banco e validado no backend.

### Usuários
| Método | Rota                | Descrição                                 |
| ------ | ------------------- | ----------------------------------------- |
| GET    | /users/:id          | Busca perfil de usuário (protegida)      |
| PATCH  | /users/:id          | Atualiza perfil do usuário (autenticado) |
| GET    | /users/:id/artworks | Lista artes de um usuário                |

### Artes
| Método | Rota               | Descrição                                  |
| ------ | ------------------ | ------------------------------------------ |
| GET    | /artworks          | Lista todas as artes (com filtros)        |
| POST   | /artworks          | Cria uma arte (upload + form-data, auth) |
| GET    | /artworks/:id      | Busca arte por ID                         |
| PATCH  | /artworks/:id      | Atualiza arte (autenticado e dono)        |
| DELETE | /artworks/:id      | Deleta a arte (autenticado e dono)        |
| POST   | /artworks/:id/like | Curte uma arte                            |
| GET    | /top-artworks      | Lista artes mais curtidas                 |
| GET    | /gallery           | Galeria web simples (HTML)                |

---

## 🔒 Protegendo rotas

Para acessar rotas protegidas, envie o JWT no header:

```
Authorization: Bearer SEU_TOKEN
```

---

#### Feito com 💙 originalmente em Go por @diozhn e agora refatorado para TypeScript