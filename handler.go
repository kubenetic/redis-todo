package main

import (
    "html/template"
    "log"
    "net/http"
)

var (
	indexTpl *template.Template
	cache *TodoCache
)

type Payload struct {
    Todos []Todo
    Error string
}

func init() {
    indexTpl = template.Must(template.ParseFiles("templates/index.gohtml"))
    cache = CreateTodoCache()
}

func handleIndex() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        todos, err := cache.GetAll()
        if err != nil {
            log.Printf("Todo lista lekérdezése sikertelen. %v\n", err)
            http.Error(w, "Index oldal kiszolgálása sikertelen", http.StatusInternalServerError)
            return
        }
        
        p := Payload{
            Todos: todos,
        } 
        if err := indexTpl.ExecuteTemplate(w, "index.gohtml", p); err != nil {
            log.Printf("Index oldal kiszolgálása sikertelen. %v\n", err)
            http.Error(w, "Index oldal kiszolgálása sikertelen", http.StatusInternalServerError)
            return
        }
    }
}

func handleAdd() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            log.Printf("todo elem hozzáadása sikertelen. %v\n", err)
            p := Payload{
                Error: "Adatok feldolgozása sikertelen",
            }

            w.WriteHeader(http.StatusBadRequest)
            _ = indexTpl.ExecuteTemplate(w, "index.gohtml", p)
            return
        }

        name := r.Form.Get("name")
        desc := r.Form.Get("description")

        t := NewTodo(name, desc)
        if err := cache.Set(t.ID, *t); err != nil {
            log.Printf("todo elem hozzáadása sikertelen. %v\n", err)
            p := Payload{
                Error: "Elem hozzáadása sikertelen",
            }

            w.WriteHeader(http.StatusInternalServerError)
            _ = indexTpl.ExecuteTemplate(w, "index.gohtml", p)
        } else {
            http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        }
    }
}

func handleDone() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            log.Printf("todo állapotának beállítása sikertelen. %v\n", err)
            p := Payload{
                Error: "Adatok feldolgozása sikertelen",
            }

            w.WriteHeader(http.StatusBadRequest)
            _ = indexTpl.ExecuteTemplate(w, "index.gohtml", p)
            return
        }

        id := r.Form.Get("ID")
        t, err := cache.Get(id)
        if err != nil {
            log.Printf("todo állapotának beállítása sikertelen. %v\n", err)
            p := Payload{
                Error: "Adatok feldolgozása sikertelen",
            }

            w.WriteHeader(http.StatusBadRequest)
            _ = indexTpl.ExecuteTemplate(w, "index.gohtml", p)
            return
        }

        t.Done = true
        if err := cache.Set(id, *t); err != nil {
            log.Printf("todo állapotának beállítása sikertelen. %v\n", err)
            p := Payload{
                Error: "Elem módosítása sikertelen",
            }

            w.WriteHeader(http.StatusInternalServerError)
            _ = indexTpl.ExecuteTemplate(w, "index.gohtml", p)
            return
        }

        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    }
}

func handleDel() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            log.Printf("todo állapotának beállítása sikertelen. %v\n", err)
            p := Payload{
                Error: "Adatok feldolgozása sikertelen",
            }

            w.WriteHeader(http.StatusBadRequest)
            _ = indexTpl.ExecuteTemplate(w, "index.gohtml", p)
            return
        }

        id := r.Form.Get("ID")
        if err := cache.Del(id); err != nil {
            log.Printf("todo állapotának beállítása sikertelen. %v\n", err)
            p := Payload{
                Error: "Elem törlése sikertelen",
            }

            w.WriteHeader(http.StatusInternalServerError)
            _ = indexTpl.ExecuteTemplate(w, "index.gohtml", p)
            return
        }

        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    }
}