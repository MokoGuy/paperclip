# PRD — `paperCLIp` : CLI d'exploration Paperless-NGX

## Contexte & motivation

Aujourd'hui, l'interaction avec Paperless passe exclusivement par Claude via MCP. Ça marche, mais avec 3 problèmes récurrents identifiés dans 23 sessions analysées :

1. **Le search full-text rate souvent** — jusqu'à 5 reformulations pour trouver un document. Le moteur de recherche Django est limité.
2. **Pas de filtres composables** — on peut chercher par texte OU filtrer par type/correspondent/date, jamais les deux ensemble en une commande.
3. **L'extraction batch est fastidieuse** — lire le contenu de N documents = N+1 appels MCP séquentiels.

**Pourquoi une CLI plutôt que le MCP ?** Shell pipelines, scripts, usage hors-Claude, et surtout : un cache local avec fuzzy-match côté client qui résout le problème de search que le MCP ne peut pas résoudre.

## Non-goals (MVP)

- Pas d'écriture : pas de `create`, `update`, `delete`, `upload`
- Pas de custom fields
- Pas de gestion de permissions/utilisateurs
- Pas de TUI interactive (possible en v2, pas en v1)
- Pas de download de fichiers PDF (consultation contenu texte uniquement)

## Cas d'usage MVP

### 1. Recherche composable avec fuzzy-match local

Le cas dominant (100% des sessions). Combiner texte libre + filtres structurés en une commande.

```bash
# Recherche full-text simple
paperclip search "invoice"

# Recherche composable : texte + type + correspondent + période
paperclip search "payslip" --type "Payslip" --from acme-corp --year 2025

# Fuzzy-match sur les noms (résout les typos et abréviations)
paperclip search --from amzn --type invoice "keyboard"
# → résout "amzn" → "Amazon", matche les factures contenant "keyboard"

# Derniers documents ajoutés
paperclip search --recent 10
```

**Mécanisme** : la CLI maintient un **cache local SQLite** de la taxonomie (tags, types, correspondents) + métadonnées des documents (titre, date, IDs). La recherche est d'abord locale (fuzzy-match sur titres + résolution des noms de filtres), puis complétée par l'API si besoin.

### 2. Extraction batch de contenu

Lire le contenu texte de plusieurs documents, pipe-friendly.

```bash
# Contenu d'un document
paperclip content 42

# Contenu de plusieurs documents
paperclip content 42 43 45 50

# Pipeline : rechercher puis extraire
paperclip search "payslip" --from acme-corp --year 2025 --ids-only | xargs paperclip content

# Extraire et grep
paperclip content 42 43 | grep "net salary"
```

### 3. Exploration taxonomique

Orientation rapide : qu'est-ce qui existe dans l'instance ?

```bash
# Lister les tags avec compteurs
paperclip tags
# Finance (425)  Employment (125)  Equipment (91)  Housing (72)  ...

# Lister les types de documents
paperclip types
# Bank Statement (223)  Invoice (130)  Payslip (83)  ...

# Lister les correspondants (top 20 par nombre de docs)
paperclip correspondents
# My Bank (217)  Acme Corp (59)  BigCo (47)  ...

# Filtrer
paperclip correspondents --filter "ban"
# My Bank (217)  Other Bank (33)
```

## Décision d'architecture : cache local SQLite

**Option retenue : (B) Cache local + fuzzy-match client.**

Raison : c'est le seul moyen de résoudre le problème #1 (search failures). Un thin wrapper REST reproduirait exactement les mêmes échecs que le MCP.

| Composant | Détail |
|-----------|--------|
| Cache | SQLite dans `~/.config/paperclip/cache.db` |
| Contenu | Métadonnées documents (titre, date, correspondent, type, tags) + taxonomie complète |
| Sync | `paperclip sync` — lazy sync avec seuil 24h |
| Fuzzy | `sahilm/fuzzy` — Levenshtein sur titres + noms correspondants/types |
| Fallback | Si le cache est vide → sync automatique au premier appel |

Le contenu texte des documents n'est **pas** caché (trop volumineux, sensible). Seul `content` reste un appel API live.

## Config & auth

```toml
# ~/.config/paperclip/config.toml
url = "https://your-paperless-instance.example.com"
token = "your-api-token-here"
```

Token API Paperless en clair dans le fichier (chmod 600 enforced par la CLI).

## Stack technique

- **Langage** : Go
- **CLI framework** : `cobra`
- **Cache** : `modernc.org/sqlite` (pure Go, pas de CGO)
- **Fuzzy** : `sahilm/fuzzy`
- **HTTP** : `net/http` (stdlib)
- **Output** : `charmbracelet/lipgloss` pour le formatage terminal
- **DB codegen** : SQLC pour les requêtes type-safe
- **Build** : binaire statique Linux

## Distribution

- **Hébergement** : GitHub
- **Artifact** : binaire exécutable Linux unique, sans dépendances
- **Installation** : download depuis les releases GitHub, ou `go install`

## Format de sortie : dual humain / agent LLM

La CLI est conçue pour être utilisée autant par un humain en terminal que par un agent LLM typé.

| Contexte | Format | Déclencheur |
|----------|--------|-------------|
| Terminal (TTY) | Table colorée, lisible | Défaut si stdout est un TTY |
| Pipe / agent LLM | JSON structuré | Auto si stdout n'est pas un TTY |
| Forcer JSON en terminal | JSON | `--json` |

Le JSON expose un **schéma stable** avec systématiquement les IDs (pour chaîner `search → content`) et les URLs web (`{base_url}/documents/{id}/`). L'agent n'a pas besoin de penser à `--json`, le pipe le déclenche automatiquement.

## Stratégie de sync : lazy avec seuil

| Situation | Comportement |
|-----------|-------------|
| Premier appel, pas de cache | Sync complète automatique |
| Cache < 24h | Utilise le cache local, instantané |
| Cache > 24h | Re-sync silencieuse au prochain appel |
| `paperclip sync` | Force une sync manuelle |
| `paperclip search --no-cache` | Bypass le cache, requête API directe (données fraîches garanties) |

Le contenu texte des documents n'est **pas** caché (trop volumineux, sensible). `content` reste toujours un appel API live.
