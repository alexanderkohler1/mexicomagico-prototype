package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Product represents a Mexican food product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
	InStock     bool    `json:"in_stock"`
}

// User represents a website user/customer
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Location string `json:"location"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

// Recipe represents a Mexican recipe
type Recipe struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	PrepTime    string   `json:"prep_time"`
	Difficulty  string   `json:"difficulty"`
	Image       string   `json:"image"`
	Ingredients []string `json:"ingredients"`
}

// PageData contains data for template rendering
type PageData struct {
	Title      string     `json:"title"`
	Products   []Product  `json:"products"`
	Recipes    []Recipe   `json:"recipes"`
	User       *User      `json:"user,omitempty"`
	Cart       []CartItem `json:"cart"`
	StoreHours string     `json:"store_hours"`
	ShowAppCTA bool       `json:"show_app_cta"`
}

// Mock data for demonstration - 5 Hauptkategorien
var products = []Product{
	// Getränke
	{1, "Mezcal Artesanal", "Handgefertigter Premium Mezcal aus Oaxaca", 64.90, "/static/images/products/mezcal-artesanal.jpg", "Getränke", true},
	{2, "Tequila Blanco", "100% Agave Tequila, kristallklar und rein", 49.90, "https://via.placeholder.com/300x200/E8F5E8/2D5016?text=Tequila+Blanco", "Getränke", true},
	{3, "Horchata de Arroz", "Traditionelles Reisgetränk mit Zimt", 3.50, "https://via.placeholder.com/300x200/E8F5E8/2D5016?text=Horchata", "Getränke", true},
	{4, "Pulque Tradicional", "Fermentierter Agavensaft, mild alkoholisch", 8.90, "https://via.placeholder.com/300x200/E8F5E8/2D5016?text=Pulque", "Getränke", true},
	
	// Basis & Grundzutaten
	{5, "Corn Tortillas (12 Stk)", "Traditionelle Maistortillas, handgemacht", 4.50, "https://via.placeholder.com/300x200/F5E8D0/8B4513?text=Tortillas", "Basis & Grundzutaten", true},
	{6, "Masa Harina", "Spezialmehl für Tortillas und Tamales", 6.90, "https://via.placeholder.com/300x200/F5E8D0/8B4513?text=Masa+Harina", "Basis & Grundzutaten", true},
	{7, "Frijoles Negros", "Schwarze Bohnen in Dose, Bio-Qualität", 2.90, "https://via.placeholder.com/300x200/F5E8D0/8B4513?text=Frijoles", "Basis & Grundzutaten", true},
	{8, "Nopal Kaktusblätter", "Eingelegte Kaktusblätter, reich an Nährstoffen", 8.90, "https://via.placeholder.com/300x200/F5E8D0/8B4513?text=Nopal", "Basis & Grundzutaten", true},
	{9, "Chayote", "Exotisches Kürbisgewächs aus Mexiko", 3.50, "https://via.placeholder.com/300x200/F5E8D0/8B4513?text=Chayote", "Basis & Grundzutaten", true},
	
	// Salsas & Saucen
	{10, "Salsa Verde", "Grüne Tomatillo-Salsa, mittelscharf", 6.90, "https://via.placeholder.com/300x200/FFE8E8/C8102E?text=Salsa+Verde", "Salsas & Saucen", true},
	{11, "Salsa Roja", "Rote Tomaten-Salsa, scharf", 6.90, "https://via.placeholder.com/300x200/FFE8E8/C8102E?text=Salsa+Roja", "Salsas & Saucen", true},
	{12, "Salsa Macha", "Nuss-Chili-Öl, premium handgemacht", 12.90, "https://via.placeholder.com/300x200/FFE8E8/C8102E?text=Salsa+Macha", "Salsas & Saucen", true},
	{13, "Mole Poblano", "Komplexe Schokoladen-Chili-Sauce", 18.90, "https://via.placeholder.com/300x200/FFE8E8/C8102E?text=Mole+Poblano", "Salsas & Saucen", true},
	{14, "Chipotle en Adobo", "Geräucherte Jalapeños in Adobo-Sauce", 4.90, "https://via.placeholder.com/300x200/FFE8E8/C8102E?text=Chipotle", "Salsas & Saucen", true},
	
	// Gewürze & Chilis
	{15, "Dried Chili Mix", "Verschiedene getrocknete Chilis aus Mexiko", 12.90, "https://via.placeholder.com/300x200/E8FFE8/006341?text=Chili+Mix", "Gewürze & Chilis", false},
	{16, "Epazote", "Traditionelles mexikanisches Kraut", 3.90, "https://via.placeholder.com/300x200/E8FFE8/006341?text=Epazote", "Gewürze & Chilis", true},
	{17, "Achiote Paste", "Rote Gewürzmischung aus Yucatán", 7.90, "https://via.placeholder.com/300x200/E8FFE8/006341?text=Achiote", "Gewürze & Chilis", true},
	{18, "Tajín Clásico", "Chili-Limetten-Gewürz, der Klassiker", 4.50, "https://via.placeholder.com/300x200/E8FFE8/006341?text=Tajin", "Gewürze & Chilis", true},
	{19, "Flor de Sal", "Meersalz aus Guerrero", 8.90, "https://via.placeholder.com/300x200/E8FFE8/006341?text=Flor+de+Sal", "Gewürze & Chilis", true},
	
	// Süßes & Snacks
	{20, "Chocolate Abuelita", "Mexikanische Trinkschokolade mit Zimt", 4.90, "https://via.placeholder.com/300x200/FFF8E8/FF6B35?text=Chocolate", "Süßes & Snacks", true},
	{21, "Dulce de Leche", "Karamellcreme nach mexikanischer Art", 6.50, "https://via.placeholder.com/300x200/FFF8E8/FF6B35?text=Dulce+de+Leche", "Süßes & Snacks", true},
	{22, "Tamarindo Candy", "Süß-saure Tamarinden-Süßigkeit", 2.90, "https://via.placeholder.com/300x200/FFF8E8/FF6B35?text=Tamarindo", "Süßes & Snacks", true},
	{23, "Churros Mix", "Fertigmischung für klassische Churros", 5.90, "https://via.placeholder.com/300x200/FFF8E8/FF6B35?text=Churros+Mix", "Süßes & Snacks", true},
	{24, "Mazapán Mexicano", "Traditionelle Erdnuss-Süßigkeit", 3.50, "https://via.placeholder.com/300x200/FFF8E8/FF6B35?text=Mazapan", "Süßes & Snacks", true},
}

var recipes = []Recipe{
	{1, "Guacamole Tradicional", "Cremiger Avocado-Dip mit Limette und Koriander", "15 Min", "Einfach", "https://via.placeholder.com/350x250/FFF8E1/8B4513?text=Guacamole", []string{"3 reife Avocados", "2 Limetten", "1 kleine Zwiebel", "2 Tomaten", "Koriander", "Salz"}},
	{2, "Tacos al Pastor", "Würzige Schweinefleisch-Tacos mit Ananas", "45 Min", "Mittel", "https://via.placeholder.com/350x250/FFF8E1/8B4513?text=Tacos+al+Pastor", []string{"500g Schweinefleisch", "Corn Tortillas", "Ananas", "Zwiebeln", "Salsa Verde", "Koriander"}},
	{3, "Quesadillas con Nopal", "Käse-Quesadillas mit Kaktusblättern", "20 Min", "Einfach", "https://via.placeholder.com/350x250/FFF8E1/8B4513?text=Quesadillas", []string{"Tortillas", "Käse", "Nopal Kaktusblätter", "Zwiebeln", "Gewürze"}},
	{4, "Chiles Rellenos", "Gefüllte Chilis mit Käse", "60 Min", "Schwer", "https://via.placeholder.com/350x250/FFF8E1/8B4513?text=Chiles+Rellenos", []string{"Poblano Chilis", "Käse", "Eier", "Mehl", "Dried Chili Mix"}},
	{5, "Mole Poblano", "Komplexe Schokoladen-Chili-Sauce", "2 Std", "Schwer", "https://via.placeholder.com/350x250/FFF8E1/8B4513?text=Mole+Poblano", []string{"Verschiedene Chilis", "Schokolade", "Nüsse", "Gewürze", "Hühnerbrühe"}},
	{6, "Horchata Casera", "Selbstgemachte Reis-Zimt-Milch", "30 Min", "Einfach", "https://via.placeholder.com/350x250/FFF8E1/8B4513?text=Horchata", []string{"Reis", "Zimt", "Milch", "Zucker", "Vanille"}},
}

var currentUser = &User{1, "Maria Schmidt", "maria@example.com", "Stuttgart"}
var cart = []CartItem{}

func Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	// Route handling
	switch {
	case path == "/" || path == "/index.html":
		homeHandler(w, r)
	case strings.HasPrefix(path, "/products"):
		productsHandler(w, r)
	case strings.HasPrefix(path, "/recipes"):
		recipesHandler(w, r)
	case strings.HasPrefix(path, "/product/"):
		productDetailHandler(w, r)
	case path == "/cart":
		cartHandler(w, r)
	case path == "/checkout":
		checkoutHandler(w, r)
	case path == "/order-confirmation":
		orderConfirmationHandler(w, r)
	case path == "/store-info":
		storeInfoHandler(w, r)
	case path == "/app-download":
		appDownloadHandler(w, r)
	case path == "/about":
		aboutHandler(w, r)
	case strings.HasPrefix(path, "/api/"):
		apiHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Show one product from each category on homepage
	featuredProducts := []Product{}
	categoriesSeen := map[string]bool{}
	for _, product := range products {
		if !categoriesSeen[product.Category] && len(featuredProducts) < 5 {
			featuredProducts = append(featuredProducts, product)
			categoriesSeen[product.Category] = true
		}
	}

	data := PageData{
		Title:      "México Mágico - Authentische mexikanische Lebensmittel",
		Products:   featuredProducts, // Show one product per category
		Recipes:    recipes[:3],      // Show first 3 recipes on homepage
		User:       currentUser,
		Cart:       cart,
		StoreHours: "Mo-Fr: 10:00-18:00, Sa: 10:00-16:00",
		ShowAppCTA: true,
	}

	renderTemplate(w, "home", data)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	filteredProducts := products
	if category != "" {
		filteredProducts = []Product{}
		for _, p := range products {
			if p.Category == category {
				filteredProducts = append(filteredProducts, p)
			}
		}
	}

	data := PageData{
		Title:    "Alle Produkte - México Mágico",
		Products: filteredProducts, // Show all products on products page
		User:     currentUser,
		Cart:     cart,
	}

	renderTemplate(w, "products", data)
}

func productDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var product *Product
	for _, p := range products {
		if p.ID == id {
			product = &p
			break
		}
	}

	if product == nil {
		http.NotFound(w, r)
		return
	}

	data := PageData{
		Title:    product.Name + " - México Mágico",
		Products: []Product{*product},
		User:     currentUser,
		Cart:     cart,
	}

	renderTemplate(w, "product-detail", data)
}

func cartHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Warenkorb - México Mágico",
		User:  currentUser,
		Cart:  cart,
	}

	renderTemplate(w, "cart", data)
}

func storeInfoHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Laden München - México Mágico",
		User:       currentUser,
		Cart:       cart,
		StoreHours: "Mo-Fr: 10:00-18:00, Sa: 10:00-16:00",
	}

	renderTemplate(w, "store-info", data)
}

func appDownloadHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "México Mágico App - Jetzt herunterladen",
		User:  currentUser,
		Cart:  cart,
	}

	renderTemplate(w, "app-download", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Unsere Geschichte - México Mágico",
		User:  currentUser,
		Cart:  cart,
	}

	renderTemplate(w, "about", data)
}

func recipesHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "Authentische Rezepte - México Mágico",
		Recipes: recipes,
		User:    currentUser,
		Cart:    cart,
	}

	renderTemplate(w, "recipes", data)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Sichere Kasse - México Mágico",
		Cart:  cart,
		User:  currentUser,
	}
	
	renderTemplate(w, "checkout", data)
}

func orderConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Bestellung bestätigt - México Mágico",
		User:  currentUser,
	}
	
	renderTemplate(w, "order-confirmation", data)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api")
	
	switch path {
	case "/products":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	case "/cart/add":
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Simple mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":    true,
			"cart_count": len(cart) + 1,
		})
	case "/user":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentUser)
	default:
		http.NotFound(w, r)
	}
}

func renderTemplate(w http.ResponseWriter, templateName string, data PageData) {
	// Embedded template content - simplified version for Vercel
	baseTemplate := `<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
        <div class="container">
            <a class="navbar-brand" href="/"><i class="fas fa-pepper-hot me-2"></i>México Mágico</a>
            <div class="navbar-nav">
                <a class="nav-link" href="/">Home</a>
                <a class="nav-link" href="/products">Produkte</a>
                <a class="nav-link" href="/recipes">Rezepte</a>
                <a class="nav-link" href="/cart">Warenkorb ({{len .Cart}})</a>
            </div>
        </div>
    </nav>
    <main class="container mt-4">
        <h1>{{.Title}}</h1>
        {{if eq . "home"}}
        <p>Willkommen bei México Mágico - Ihre Quelle für authentische mexikanische Lebensmittel!</p>
        {{else if eq . "products"}}
        <div class="row">
            {{range .Products}}
            <div class="col-md-4 mb-3">
                <div class="card">
                    <div class="card-body">
                        <h5>{{.Name}}</h5>
                        <p>{{.Description}}</p>
                        <p class="h6">{{printf "%.2f" .Price}}€</p>
                    </div>
                </div>
            </div>
            {{end}}
        </div>
        {{else}}
        <p>Seite wird geladen...</p>
        {{end}}
    </main>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`

	// Parse and execute template
	tmpl, err := template.New("base").Parse(baseTemplate)
	if err != nil {
		log.Printf("Template parse error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}