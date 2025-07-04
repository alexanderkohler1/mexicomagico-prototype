package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Product represents a Mexican food product
type Product struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Image          string  `json:"image"`
	Category       string  `json:"category"`
	InStock        bool    `json:"in_stock"`
	RequiresAgeVerification bool `json:"requires_age_verification"`
	AvailabilityStatus string `json:"availability_status"`
	IsTopSeller    bool    `json:"is_top_seller"`
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
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	PrepTime        string   `json:"prep_time"`
	Difficulty      string   `json:"difficulty"`
	Image           string   `json:"image"`
	Ingredients     []string `json:"ingredients"`
	Category        string   `json:"category"`
	IsVegan         bool     `json:"is_vegan"`
	Instructions    []string `json:"instructions"`
	AboutRecipe     string   `json:"about_recipe"`
	InStoreIngredients []string `json:"in_store_ingredients"`
	AdditionalIngredients []string `json:"additional_ingredients"`
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

var products = []Product{
	// Basis & Grundzutaten
	{1, "KAKTUSBLÄTTER NOPALES STREIFEN IN SALZLAKE 300g", "Eingelegte Kaktusblätter-Streifen, reich an Nährstoffen", 4.25, "/static/images/products/KAKTUSBLÄTTER NOPALES STREIFEN IN SALZLAKE 300g.jpg", "Basis & Grundzutaten", true, false, "available", true},
	{2, "GANZE NOPALES KAKTUSBLÄTTER IM BEUTEL VON AZTECA 1kg", "Ganze Kaktusblätter für authentische mexikanische Gerichte", 8.35, "/static/images/products/GANZE NOPALES KAKTUSBLÄTTER IM BEUTEL VON AZTECA 1kg.png", "Basis & Grundzutaten", true, false, "available", false},
	{3, "Tamal mit Chili Poblano und Käse (Rajas con queso) 300 g (3 St.)", "Traditionelle Tamales mit Poblano-Chili und Käse", 11.49, "/static/images/products/Tamal mit Chili Poblano und Käse (Rajas con queso) 300 g (3 St.).png", "Basis & Grundzutaten", true, false, "available", true},
	{4, "Chorizoqueso 500g", "Würzige Chorizo-Käse-Mischung", 13.50, "/static/images/products/Chorizoqueso 500g.png", "Basis & Grundzutaten", true, false, "available", false},
	{5, "Cochinita Pibil (servierfertig) 250g", "Servierfertige Cochinita Pibil aus Yucatán", 9.90, "/static/images/products/Cochinita Pibil (servierfertig) 250g.png", "Basis & Grundzutaten", true, false, "available", false},
	{16, "Totopos Nixtamal Tortilla-Chips gesalzen 250g", "Authentische Tortilla-Chips aus Nixtamal-Mais", 6.40, "/static/images/products/Totopos Nixtamal Tortilla-Chips gesalzen 250g Authentische Tortilla-Chips aus Nixtamal-Mais.png", "Basis & Grundzutaten", false, false, "coming_soon", false},
	{17, "frische Maistortillas aus Nixtamal 15cm ca. 10 Stk. Glutenfrei", "Traditionelle Maistortillas, handgemacht", 4.69, "/static/images/products/frische Maistortillas aus Nixtamal 15cm ca. 10 Stk. Glutenfrei.png", "Basis & Grundzutaten", true, false, "available", false},
	{18, "MAIS TOSTADAS VON CHARRAS 325 g", "Knusprige Mais-Tostadas für authentische Gerichte", 6.59, "/static/images/products/MAIS TOSTADAS VON CHARRAS 325 g.png", "Basis & Grundzutaten", true, false, "available", false},

	// Getränke
	{10, "BOING MANGO SAFT 354 ML", "Natürlicher Mango-Saft aus Mexiko", 3.50, "/static/images/products/BOING MANGO SAFT 354 ML.png", "Getränke", true, false, "available", true},
	{7, "400 Conejos Mezcal Joven 38% Vol. 0.7 l", "Authentischer Mezcal Joven aus Oaxaca", 49.90, "/static/images/products/400 Conejos Mezcal Joven 38% Vol. 0.7 l.png", "Getränke", true, true, "available", false},
	{8, "Mezcal Don Ramón Joven 100% Agave Salmiana 40% Vol. 0,7l", "Premium Mezcal aus 100% Agave Salmiana", 52.00, "/static/images/products/Mezcal Don Ramón Joven 100% Agave Salmiana 40% Vol. 0,7l.png", "Getränke", false, true, "coming_soon", false},
	{9, "MODELO ESPECIAL HELLES BIER 355 ml 4.5% Vol. Alk.", "Authentisches mexikanisches Lagerbier", 3.69, "/static/images/products/MODELO ESPECIAL HELLES BIER 355 ml 4.5% Vol. Alk..png", "Getränke", true, true, "available", false},
	{11, "HORCHATA KONZENTRAT VON MEXQUISITA 700 ML", "Konzentrat für traditionelle Horchata", 7.60, "/static/images/products/HORCHATA KONZENTRAT VON MEXQUISITA 700 ML.png", "Getränke", false, false, "coming_soon", false},

	// Salsas & Saucen
	{12, "SCHARFE SOßE LA BOTANERA CLASICA 1L", "Klassische scharfe Sauce aus Mexiko", 5.99, "/static/images/products/SCHARFE SOße LA BOTANERA CLASICA 1L.png", "Salsas & Saucen", true, false, "available", true},
	{13, "HABANERO EXTRASCHARFE SOßE VON EL YUCATECO 120ml", "Extrascharfe Habanero-Sauce", 4.49, "/static/images/products/HABANERO EXTRASCHARFE SOße VON EL YUCATECO 120ml.png", "Salsas & Saucen", true, false, "available", false},
	{14, "Schwarze Habanero Chili Soße El Yucateco 120ml", "Schwarze Habanero-Sauce mit einzigartigem Geschmack", 4.49, "/static/images/products/Schwarze Habanero Chili Soße El Yucateco 120ml.png", "Salsas & Saucen", true, false, "available", false},
	{15, "VALENTINA EXTRA SCHARFE SAUCE 370ml", "Mexikos beliebteste scharfe Sauce", 3.69, "/static/images/products/VALENTINA EXTRA SCHARFE SAUCE 370ml.png", "Salsas & Saucen", true, false, "available", false},

	// Gewürze & Chilis
	{6, "CHILI ARBOL VON XATZE 75gr", "Getrocknete Chili de Árbol für intensive Schärfe", 5.49, "/static/images/products/CHILI ARBOL VON XATZE.png", "Gewürze & Chilis", true, false, "available", false},
	{19, "TAJIN CHILI LIMETTEN PULVER 142gr", "Das Original Tajín Gewürz mit Chili und Limette", 3.79, "/static/images/products/TAJIN CHILI LIMETTEN PULVER 142gr.png", "Gewürze & Chilis", true, false, "available", false},
	{20, "GEWÜRZ FÜR FAJITA - COCHINITA PIBIL 142 g", "Authentische Gewürzmischung für Cochinita Pibil", 5.50, "/static/images/products/GEWÜRZ FÜR FAJITA - COCHINITA PIBIL 142 g.png", "Gewürze & Chilis", true, false, "available", false},
	{21, "MEXIKANISCHER OREGANO (GEWÜRZ) VON XATZE 20gr", "Echter mexikanischer Oregano", 3.99, "https://mexicomagico.de/cdn/shop/files/Oregano-para-pozole-20147-_1.jpg?v=1683208957&width=400", "Gewürze & Chilis", true, false, "available", false},

	// Süßes & Snacks
	{22, "Vero Mango Lollis mit Chili 16g", "Süß-scharfe Mango-Lollis mit Chili", 1.20, "/static/images/products/Vero Mango Lollis mit Chili 16g.png", "Süßes & Snacks", true, false, "available", false},
	{23, "Pulparindo Süßigkeit von De La Rosa 280g (20 St.)", "Süß-saure Tamarinden-Süßigkeit", 7.90, "https://mexicomagico.de/cdn/shop/files/FotosProductos_16.png?v=1683208957&width=400", "Süßes & Snacks", true, false, "available", false},
	{24, "Pelon Pelo Rico Tamarinde mit Chili 35g", "Tamarinden-Süßigkeit mit Chili", 1.49, "/static/images/products/Pelon Pelo Rico Tamarinde mit Chili 35g.png", "Süßes & Snacks", true, false, "available", false},
	{25, "Pulparindos CHAMOY 20 St. (280 g)", "Chamoy-Süßigkeiten von De La Rosa", 7.90, "/static/images/products/Pulparindos CHAMOY 20 St. (280 g).png", "Süßes & Snacks", true, false, "available", false},
}

var recipes = []Recipe{
	{1, "Nopales- und Avocado-Salat", "Frischer Salat mit Kaktusblättern und Avocado", "15 Min", "Einfach", "/static/images/recipes/NopalesUndAvocadoSalat.png", []string{"500g Geschnittene Kaktusblätter", "1 Jalapeño (optional)", "Glutenfreie Maistortillas", "2 reife Avocados", "1 große Tomate", "1/4 rote Zwiebel", "Saft von 1 Zitrone", "2 EL frisch gehackter Koriander", "Salz und Pfeffer", "3 EL Feta-Käse"}, "Hauptgerichte", true, []string{"1. Vorbereitung der Nopales: Die Nopales unter fließendem Wasser waschen, um Schmutz zu entfernen. Die Stacheln der Nopales vorsichtig mit einem Messer oder einem Sparschäler entfernen. Die Nopales in Streifen oder Würfel schneiden, je nach Vorliebe. In einer Pfanne oder auf einem Grill etwas Olivenöl bei mittlerer bis hoher Hitze erhitzen. Die Nopales in der heißen Pfanne etwa 5-7 Minuten auf jeder Seite grillen, bis sie zart und leicht gebräunt sind. Vom Herd nehmen und in eine große Schüssel geben.", "2. Salatzubereitung: Die Avocados in Würfel und die Tomate in kleine Stücke schneiden. Die rote Zwiebel und die Jalapeño (falls verwendet) fein hacken. Die Avocados, Tomaten, rote Zwiebel und Jalapeño in die Schüssel mit den gegrillten Nopales geben. Den Zitronensaft über den Salat auspressen und den frisch gehackten Koriander hinzufügen. Alle Zutaten vorsichtig vermengen, bis sie gut kombiniert sind. Den Salat mit Salz und Pfeffer abschmecken.", "3. Servieren: Die glutenfreien Maistortillas in einer Pfanne oder auf einem Grill erwärmen, bis sie warm und leicht gebräunt sind. Den Nopales- und Avocado-Salat auf einem Teller servieren und mit den glutenfreien Maistortillas in Form von Tacos begleiten. Zum Servieren den geriebenen Feta-Käse über den Salat streuen."}, "Der Nopal ist ein äußerst nahrhaftes Lebensmittel, das eine breite Palette an gesundheitlichen Vorteilen bietet. Er ist reich an Ballaststoffen, was eine gesunde Verdauung fördert und das Gewicht kontrolliert. Es hilft, den Blutzuckerspiegel zu kontrollieren und das LDL-Cholesterin zu senken. Darüber hinaus ist er eine Quelle von Antioxidantien, die freie Radikale bekämpfen, und liefert Feuchtigkeit aufgrund seines hohen Wassergehalts. Er ist reich an essentiellen Vitaminen und Mineralstoffen wie Vitamin A, C, K, Calcium, Kalium und Eisen. Er hat auch entzündungshemmende Eigenschaften, was ihn vorteilhaft für Bedingungen wie Arthritis macht. Zusammenfassend verbessert der regelmäßige Verzehr von Nopal die Verdauungs-, Herz-Kreislauf- und Immunfunktion und hilft beim Gewichtsmanagement sowie der Feuchtigkeitsversorgung.", []string{"500g Geschnittene Kaktusblätter", "1 Jalapeño (optional, für eine pikante Note)", "Glutenfreie Maistortillas (je nach Anzahl der Portionen)"}, []string{"2 reife Avocados", "1 große Tomate", "1/4 rote Zwiebel", "Saft von 1 Zitrone", "2 Esslöffel frisch gehackter Koriander", "Salz und Pfeffer nach Geschmack", "3 Esslöffel Feta-Käse"}},
	{2, "Hähnchen mit Mole", "Zartes Hähnchen mit komplexer Mole-Sauce", "2 Std", "Schwer", "/static/images/recipes/ZartesHähnchenMitKomplexerMoleSauce.png", []string{"1 ganzes Hähnchen", "Mole Poblano", "Verschiedene Chilis", "Schokolade", "Nüsse", "Gewürze"}, "Hauptgerichte", false, []string{"Hähnchen waschen und trockentupfen", "Mole-Sauce nach Anleitung zubereiten", "Hähnchen garen und mit Sauce servieren"}, "Ein klassisches mexikanisches Gericht aus Puebla mit komplexer, schokoladiger Sauce.", []string{"1 ganzes Hähnchen", "Mole Poblano Paste"}, []string{"Verschiedene Chilis", "Schokolade", "Nüsse", "Gewürze"}},
	{3, "Chili Poblano (rajas) mit Sahne", "Poblano-Chilistreifen in cremiger Sahnesauce", "30 Min", "Mittel", "/static/images/recipes/PoblanoChiliStreifenInCremigerSahnesauce.png", []string{"Poblano Chilis", "Sahne", "Zwiebeln", "Knoblauch", "Gewürze", "Öl"}, "Hauptgerichte", true, []string{"Poblanos rösten und häuten", "In Streifen schneiden", "Mit Sahne und Zwiebeln zubereiten"}, "Ein cremiges, vegetarisches Gericht mit gerösteten Poblano-Chilis.", []string{"Poblano Chilis", "Sahne"}, []string{"Zwiebeln", "Knoblauch", "Gewürze", "Öl"}},
	{4, "Mexikanische rote Pozole", "Traditionelle Hominy-Suppe mit roter Chilisauce", "3 Std", "Schwer", "/static/images/recipes/Traditionelle Hominy-Suppe mit roter Chilisauce.png", []string{"Hominy", "Schweinefleisch", "Guajillo Chilis", "Ancho Chilis", "Zwiebeln", "Knoblauch", "Oregano"}, "Hauptgerichte", false, []string{"Schweinefleisch kochen", "Chilis rösten und pürieren", "Hominy hinzufügen und köcheln lassen"}, "Eine traditionelle mexikanische Suppe, perfekt für Feiertage.", []string{"Hominy", "Schweinefleisch"}, []string{"Guajillo Chilis", "Ancho Chilis", "Zwiebeln", "Knoblauch", "Oregano"}},
	{5, "Michelada mit Clamato", "Würziges Bier-Cocktail mit Clamato-Saft", "5 Min", "Einfach", "/static/images/recipes/Würziges Bier-Cocktail mit Clamato-Saft.jpeg", []string{"Bier", "Clamato-Saft", "Limettensaft", "Worcestershire-Sauce", "Tabasco", "Salz", "Tajín"}, "Getränke", false, []string{"Glas mit Tajín salzen", "Alle Zutaten mischen", "Mit Eis servieren"}, "Der perfekte mexikanische Cocktail für warme Tage.", []string{"Bier", "Clamato-Saft"}, []string{"Limettensaft", "Worcestershire-Sauce", "Tabasco", "Salz", "Tajín"}},
	{6, "Flan napolitano", "Klassischer mexikanischer Karamellpudding", "4 Std", "Mittel", "/static/images/recipes/Klassischer mexikanischer Karamellpudding.png", []string{"Eier", "Kondensmilch", "Evaporierte Milch", "Zucker", "Vanille"}, "Desserts", false, []string{"Karamell zubereiten", "Pudding-Mischung vorbereiten", "Im Wasserbad backen und kühlen"}, "Ein cremiger, süßer Abschluss jeder mexikanischen Mahlzeit.", []string{"Eier", "Kondensmilch", "Evaporierte Milch"}, []string{"Zucker", "Vanille"}},
}

var currentUser = &User{1, "Maria Schmidt", "maria@example.com", "Stuttgart"}
var cart = []CartItem{}

func main() {
	// Serve static files (CSS, JS, Images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/recipes", recipesHandler)
	http.HandleFunc("/recipe/", recipeDetailHandler)
	http.HandleFunc("/product/", productDetailHandler)
	http.HandleFunc("/cart", cartHandler)
	http.HandleFunc("/checkout", checkoutHandler)
	http.HandleFunc("/order-confirmation", orderConfirmationHandler)
	http.HandleFunc("/store-info", storeInfoHandler)
	http.HandleFunc("/app-download", appDownloadHandler)
	http.HandleFunc("/about", aboutHandler)

	// API endpoints
	http.HandleFunc("/api/products", apiProductsHandler)
	http.HandleFunc("/api/cart/add", apiAddToCartHandler)
	http.HandleFunc("/api/cart/update", apiUpdateCartHandler)
	http.HandleFunc("/api/cart/remove", apiRemoveFromCartHandler)
	http.HandleFunc("/api/user", apiUserHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
	idStr := r.URL.Path[len("/product/"):]
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

// API Handlers
func apiProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func apiAddToCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Find product
	var product *Product
	for _, p := range products {
		if p.ID == reqData.ProductID {
			product = &p
			break
		}
	}

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Add to cart or update quantity
	found := false
	for i := range cart {
		if cart[i].Product.ID == reqData.ProductID {
			cart[i].Quantity += reqData.Quantity
			found = true
			break
		}
	}

	if !found {
		cart = append(cart, CartItem{
			Product:  *product,
			Quantity: reqData.Quantity,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"cart_count": len(cart),
	})
}

func apiUpdateCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update quantity in cart
	for i := range cart {
		if cart[i].Product.ID == reqData.ProductID {
			if reqData.Quantity <= 0 {
				// Remove item if quantity is 0 or less
				cart = append(cart[:i], cart[i+1:]...)
			} else {
				cart[i].Quantity = reqData.Quantity
			}
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"cart_count": len(cart),
	})
}

func apiRemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData struct {
		ProductID int `json:"product_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Remove item from cart
	for i := range cart {
		if cart[i].Product.ID == reqData.ProductID {
			cart = append(cart[:i], cart[i+1:]...)
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"cart_count": len(cart),
	})
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentUser)
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

func recipeDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL
	idStr := r.URL.Path[len("/recipe/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var recipe *Recipe
	for _, rec := range recipes {
		if rec.ID == id {
			recipe = &rec
			break
		}
	}

	if recipe == nil {
		http.NotFound(w, r)
		return
	}

	data := PageData{
		Title:   recipe.Name + " - México Mágico",
		Recipes: []Recipe{*recipe},
		User:    currentUser,
		Cart:    cart,
	}

	renderTemplate(w, "recipe-detail", data)
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

func renderTemplate(w http.ResponseWriter, templateName string, data PageData) {
	var templateFiles []string

	switch templateName {
	case "home":
		templateFiles = []string{"templates/base.html", "templates/home.html"}
	case "products":
		templateFiles = []string{"templates/base.html", "templates/products.html"}
	case "recipes":
		templateFiles = []string{"templates/base.html", "templates/recipes.html"}
	case "recipe-detail":
		templateFiles = []string{"templates/base.html", "templates/recipe-detail.html"}
	case "store-info":
		templateFiles = []string{"templates/base.html", "templates/store-info.html"}
	case "app-download":
		templateFiles = []string{"templates/base.html", "templates/app-download.html"}
	case "about":
		templateFiles = []string{"templates/base.html", "templates/about.html"}
	case "product-detail":
		templateFiles = []string{"templates/base.html", "templates/product-detail.html"}
	case "cart":
		templateFiles = []string{"templates/base.html", "templates/cart.html"}
	case "checkout":
		templateFiles = []string{"templates/base.html", "templates/checkout.html"}
	case "order-confirmation":
		templateFiles = []string{"templates/base.html", "templates/order-confirmation.html"}
	default:
		templateFiles = []string{"templates/base.html", "templates/home.html"}
	}

	// Create template with custom functions
	funcMap := template.FuncMap{
		"sub": func(a, b interface{}) float64 {
			var aFloat, bFloat float64
			
			switch v := a.(type) {
			case int:
				aFloat = float64(v)
			case float64:
				aFloat = v
			default:
				aFloat = 0
			}
			
			switch v := b.(type) {
			case int:
				bFloat = float64(v)
			case float64:
				bFloat = v
			default:
				bFloat = 0
			}
			
			return aFloat - bFloat
		},
		"add": func(a, b int) int {
			return a + b
		},
		"len": func(v interface{}) int {
			switch s := v.(type) {
			case []Product:
				return len(s)
			case []Recipe:
				return len(s)
			case []CartItem:
				return len(s)
			case []string:
				return len(s)
			default:
				return 0
			}
		},
		"mul": func(a float64, b float64) float64 {
			return a * b
		},
		"float64": func(i int) float64 {
			return float64(i)
		},
		"calculateSubtotal": func(cart []CartItem) float64 {
			var total float64
			for _, item := range cart {
				total += item.Product.Price * float64(item.Quantity)
			}
			return total
		},
		"calculateTotal": func(cart []CartItem) float64 {
			subtotal := 0.0
			for _, item := range cart {
				subtotal += item.Product.Price * float64(item.Quantity)
			}
			shipping := 0.0
			if subtotal < 57.0 {
				shipping = 4.90
			}
			return subtotal + shipping
		},
		"ge": func(a, b float64) bool {
			return a >= b
		},
	}

	tmpl, err := template.New("base").Funcs(funcMap).ParseFiles(templateFiles...)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
