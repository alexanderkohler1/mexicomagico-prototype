package handler

import (
	"encoding/json"
	"fmt"
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
	var content string
	
	switch templateName {
	case "home":
		content = renderHome(data)
	case "products":
		content = renderProducts(data)
	case "recipes":
		content = renderRecipes(data)
	case "cart":
		content = renderCart(data)
	case "product-detail":
		content = renderProductDetail(data)
	case "checkout":
		content = renderCheckout(data)
	case "order-confirmation":
		content = renderOrderConfirmation(data)
	case "store-info":
		content = renderStoreInfo(data)
	case "app-download":
		content = renderAppDownload(data)
	case "about":
		content = renderAbout(data)
	default:
		content = renderHome(data)
	}
	
	// Base template with dynamic content
	baseTemplate := `<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + data.Title + `</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <style>
        :root {
            --mexico-red: #C8102E;
            --mexico-green: #006341;
            --mexico-gold: #FFD700;
            --mexico-orange: #FF6B35;
            --mexico-cream: #FFF8E1;
            --mexico-brown: #8B4513;
        }
        body { font-family: 'Georgia', serif; background-color: var(--mexico-cream); }
        .navbar { background: linear-gradient(90deg, var(--mexico-cream) 0%, #fff 100%) !important; border-bottom: 3px solid var(--mexico-red); }
        .navbar-brand { color: var(--mexico-red) !important; font-weight: bold; font-size: 1.8rem; }
        .btn-mexico { background: linear-gradient(135deg, var(--mexico-red), var(--mexico-orange)); color: white; border: none; border-radius: 25px; }
        .card:hover { transform: translateY(-5px); transition: all 0.3s ease; }
    </style>
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-light sticky-top">
        <div class="container">
            <a class="navbar-brand" href="/"><i class="fas fa-pepper-hot me-2"></i>México Mágico</a>
            <div class="navbar-nav me-auto">
                <a class="nav-link" href="/"><i class="fas fa-home me-1"></i>Home</a>
                <a class="nav-link" href="/products"><i class="fas fa-shopping-basket me-1"></i>Produkte</a>
                <a class="nav-link" href="/recipes"><i class="fas fa-utensils me-1"></i>Rezepte</a>
                <a class="nav-link" href="/store-info"><i class="fas fa-store me-1"></i>Laden</a>
                <a class="nav-link" href="/app-download"><i class="fas fa-mobile-alt me-1"></i>App</a>
                <a class="nav-link" href="/about"><i class="fas fa-heart me-1"></i>Story</a>
            </div>
            <div class="d-flex">
                <a href="/cart" class="btn btn-outline-secondary"><i class="fas fa-shopping-cart"></i> (` + fmt.Sprintf("%d", len(data.Cart)) + `)</a>
            </div>
        </div>
    </nav>
    <main>` + content + `</main>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, baseTemplate)
}

func renderHome(data PageData) string {
	content := `
    <section style="background: linear-gradient(135deg, var(--mexico-red) 0%, var(--mexico-orange) 50%, var(--mexico-green) 100%); color: white; padding: 100px 0;">
        <div class="container text-center">
            <h1 class="display-3 mb-4">¡Bienvenidos a México Mágico!</h1>
            <p class="lead mb-4">Von Rodrigos Herzen nach Europa - authentische mexikanische Lebensmittel</p>
            <a href="/products" class="btn btn-light btn-lg me-3"><i class="fas fa-shopping-basket me-2"></i>Jetzt entdecken</a>
            <a href="/about" class="btn btn-outline-light btn-lg"><i class="fas fa-heart me-2"></i>Unsere Story</a>
        </div>
    </section>
    <section class="py-5">
        <div class="container">
            <h2 class="text-center mb-5" style="color: var(--mexico-brown);">Unsere Produktkategorien</h2>
            <div class="row">`
            
	for _, product := range data.Products {
		content += fmt.Sprintf(`
                <div class="col-md-6 col-lg-3 mb-4">
                    <div class="card h-100 border-0 shadow-sm" onclick="window.location.href='/product/%d'" style="cursor: pointer;">
                        <img src="%s" class="card-img-top" alt="%s" style="height: 200px; object-fit: cover;">
                        <div class="card-body">
                            <h5>%s</h5>
                            <p class="small">%s</p>
                            <div class="d-flex justify-content-between align-items-center">
                                <span class="h6" style="color: var(--mexico-red);">%.2f€</span>
                                <button class="btn btn-mexico btn-sm" onclick="event.stopPropagation(); addToCart(%d)">
                                    <i class="fas fa-cart-plus"></i> Hinzufügen
                                </button>
                            </div>
                        </div>
                    </div>
                </div>`, product.ID, product.Image, product.Name, product.Name, product.Description, product.Price, product.ID)
	}
	
	content += `
            </div>
            <div class="text-center mt-4">
                <a href="/products" class="btn btn-mexico btn-lg">Alle Produkte entdecken</a>
            </div>
        </div>
    </section>`
    
    return content
}

func renderProducts(data PageData) string {
	content := `
    <section style="background: linear-gradient(135deg, var(--mexico-green) 0%, var(--mexico-red) 100%); color: white; padding: 80px 0;">
        <div class="container text-center">
            <h1 class="display-4 mb-3">Alle Produkte</h1>
            <p class="lead">Entdecke die ganze Vielfalt authentischer mexikanischer Lebensmittel</p>
        </div>
    </section>
    <section class="py-5">
        <div class="container">
            <div class="row">`
            
	for _, product := range data.Products {
		inStockBadge := ""
		buttonText := `<button class="btn btn-mexico btn-sm" onclick="addToCart(` + fmt.Sprintf("%d", product.ID) + `)"><i class="fas fa-cart-plus"></i> Hinzufügen</button>`
		if !product.InStock {
			inStockBadge = `<span class="badge bg-secondary position-absolute top-0 end-0 m-2">Ausverkauft</span>`
			buttonText = `<button class="btn btn-outline-secondary btn-sm" disabled><i class="fas fa-clock"></i> Bald wieder da</button>`
		}
		
		content += fmt.Sprintf(`
            <div class="col-md-6 col-lg-4 col-xl-3 mb-4">
                <div class="card h-100 border-0 shadow-sm" onclick="window.location.href='/product/%d'" style="cursor: pointer;">
                    <div class="position-relative">
                        <img src="%s" class="card-img-top" alt="%s" style="height: 200px; object-fit: cover;">
                        <span class="badge position-absolute top-0 start-0 m-2" style="background: var(--mexico-green); color: white;">%s</span>
                        %s
                    </div>
                    <div class="card-body d-flex flex-column">
                        <h5>%s</h5>
                        <p class="small">%s</p>
                        <div class="mt-auto">
                            <div class="d-flex justify-content-between align-items-center">
                                <span class="h6" style="color: var(--mexico-red);">%.2f€</span>
                                %s
                            </div>
                        </div>
                    </div>
                </div>
            </div>`, product.ID, product.Image, product.Name, product.Category, inStockBadge, product.Name, product.Description, product.Price, buttonText)
	}
	
	content += `
            </div>
        </div>
    </section>`
    
    return content
}

func renderRecipes(data PageData) string {
	return `<section class="py-5"><div class="container"><h1>Authentische Rezepte</h1><p>Entdecken Sie unsere mexikanischen Rezepte...</p></div></section>`
}

func renderCart(data PageData) string {
	if len(data.Cart) == 0 {
		return `<section class="py-5"><div class="container text-center"><h1>Warenkorb</h1><p>Ihr Warenkorb ist leer.</p><a href="/products" class="btn btn-mexico">Jetzt einkaufen</a></div></section>`
	}
	return `<section class="py-5"><div class="container"><h1>Warenkorb</h1><p>Ihre Artikel werden hier angezeigt...</p></div></section>`
}

func renderProductDetail(data PageData) string {
	if len(data.Products) == 0 {
		return `<section class="py-5"><div class="container"><h1>Produkt nicht gefunden</h1></div></section>`
	}
	product := data.Products[0]
	return fmt.Sprintf(`
    <section class="py-5">
        <div class="container">
            <div class="row">
                <div class="col-lg-6">
                    <img src="%s" class="img-fluid rounded" alt="%s" style="width: 100%%; height: 400px; object-fit: cover;">
                </div>
                <div class="col-lg-6">
                    <h1 style="color: var(--mexico-brown);">%s</h1>
                    <p class="lead">%s</p>
                    <div class="mb-4">
                        <span class="h2" style="color: var(--mexico-red);">%.2f€</span>
                    </div>
                    <button class="btn btn-mexico btn-lg" onclick="addToCart(%d)">
                        <i class="fas fa-cart-plus me-2"></i>In den Warenkorb
                    </button>
                </div>
            </div>
        </div>
    </section>`, product.Image, product.Name, product.Name, product.Description, product.Price, product.ID)
}

func renderCheckout(data PageData) string {
	return `<section class="py-5"><div class="container"><h1>Kasse</h1><p>Checkout-Prozess...</p></div></section>`
}

func renderOrderConfirmation(data PageData) string {
	return `<section class="py-5"><div class="container text-center"><h1>Bestellung bestätigt!</h1><p>Vielen Dank für Ihre Bestellung.</p></div></section>`
}

func renderStoreInfo(data PageData) string {
	return `<section class="py-5"><div class="container"><h1>Laden München</h1><p>Maximilianstraße 42, 80538 München</p></div></section>`
}

func renderAppDownload(data PageData) string {
	return `<section class="py-5"><div class="container text-center"><h1>México Mágico App</h1><p>Jetzt herunterladen und Vorteile sichern!</p></div></section>`
}

func renderAbout(data PageData) string {
	return `<section class="py-5"><div class="container"><h1>Unsere Geschichte</h1><p>Von Rodrigo & Lisa - eine Liebesgeschichte zu Mexiko...</p></div></section>`
}