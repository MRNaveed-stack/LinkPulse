# **VS Code Copilot Debugging Strategy**

## **Phase 1: Understand the Flow (Top-Down Analysis)**

When a feature isn't working, trace the complete flow:
1. **Start at the user interaction point** (What the user does)
2. **Follow the frontend code** (Components, hooks, API calls)
3. **Trace to the backend** (Routes, handlers, services)
4. **Go down to the database** (Repositories, queries, schema)

**Example:**
```
User clicks "Create Link" button
  ↓
LinkForm component (on submit)
  ↓
useCreateLink hook (mutation)
  ↓
API client sends POST /links
  ↓
Backend LinkHandler.CreateLink
  ↓
LinkService.CreateLink
  ↓
LinkRepository.Create (database insert)
  ↓
Database schema (links table)
```

---

## **Phase 2: Identify the Likely Problem Areas (Layered Investigation)**

Investigate in this order of likelihood:

### **Layer 1: UI/Component Issues** (Most obvious)
- Missing CSS classes → No borders, invisible elements
- Wrong event handlers → Click events not firing
- Form validation errors → Schema mismatches

**Action:** Read component code, look for obvious styling/logic issues

### **Layer 2: Data Flow Issues** (Frontend ↔ Backend)
- Wrong URL construction → Requests go to wrong place
- Missing/wrong parameters → Backend can't find resource
- Type mismatches → Field names don't match

**Action:** Check API client config, examine request/response types

### **Layer 3: Backend Logic Issues** (Handler/Service)
- Missing validation
- Wrong query logic
- Authorization failures

**Action:** Read handler and service code

### **Layer 4: Database Issues** (Schema/Constraints)
- Wrong constraints → Data can't be inserted
- Missing columns → Queries fail
- Schema doesn't match code expectations

**Action:** Check migration files and table structure

---

## **Phase 3: Gather Evidence (Before Making Changes)**

**Never guess. Always read the code to understand what's happening:**

1. **Read the complete component/file** - Not just pieces
2. **Check type definitions** - Understand data shapes
3. **Trace function calls** - Follow the chain of execution
4. **Examine database schema** - See what constraints exist
5. **Look at related code** - Check similar working features

**Example Investigation:**
- ✅ Read LinkCard component
- ✅ Read the redirect handler
- ✅ Checked the database schema migration
- ✅ Read the repository methods
- ✅ Examined the service layer
- ✅ Checked how the URL was constructed

---

## **Phase 4: Form a Hypothesis**

Once you have evidence, form a specific hypothesis:

| Issue | Evidence | Hypothesis |
|-------|----------|-----------|
| Form not visible | `border-gray-300` exists but no `border` class | Missing Tailwind border utility |
| Can't create links | Created migration with `UNIQUE(slug)` | Global slug constraint blocks multi-user creation |
| Redirect logs out | URL uses `/u/{user_id}/{slug}` + backend expects `/u/{username}/{slug}` | Frontend/backend URL mismatch |
| External link redirects to login | Short URL uses `window.location.origin` (frontend port 5173) | Backend API on port 8080 is being ignored |

---

## **Phase 5: Verify the Hypothesis**

Before fixing, verify:
- ✅ Does the hypothesis explain the observed behavior?
- ✅ Are there other symptoms that support this?
- ✅ Does it make logical sense?

---

## **Phase 6: Implement the Fix (Minimal Changes)**

Fix only what's needed:
- Add the missing class
- Update the constraint
- Change the URL
- Update the function signature

**Avoid:**
- ❌ Rewriting large sections
- ❌ Changing multiple things at once
- ❌ Adding "just in case" code

---

## **Phase 7: Trace the Impact (Ripple Effect)**

After fixing one component, check what else breaks:
- Does the handler signature change? → Update all callers
- Does the database schema change? → Update frontend types
- Does the API response change? → Update frontend models

**Example:**
When changing redirect to use username:
1. ✅ Updated LinkHandler (added userRepo parameter)
2. ✅ Updated LinkService (added userID parameter)
3. ✅ Updated repository (added new method)
4. ✅ Updated main.go (pass userRepo to handler)
5. ✅ Updated LinkCard (use username instead of user_id)

---

## **Key Principles**

### **1. Read, Don't Guess**
Examine actual code instead of assuming. This catches 80% of bugs immediately.

### **2. Follow the Data**
Track what data flows through each layer. Type mismatches are common.

### **3. Check Constraints**
Database constraints, API routes, type definitions—these often hide bugs.

### **4. Test One Thing at a Time**
Fix the smallest piece first. This isolates which change solved the problem.

### **5. Look for Pattern Mismatches**
When frontend and backend don't agree on:
- URL structure
- Field names
- Data types
- Authorization

This is usually where bugs hide.

### **6. Use the Error Message as a Starting Point**
- "Link not found" → Check if the query is looking in the right place
- "Unauthorized" → Check auth middleware
- "Invalid request" → Check request validation

### **7. Build a Mental Map**
As you investigate, create a mental picture of:
- Where the data enters
- What transforms it
- Where it gets stored
- How it comes back out

---

## **Tools Used During Investigation**

| Layer | Tool | What to Look For |
|-------|------|-----------------|
| Frontend | Read component code | Props, handlers, API calls |
| Types | Read .ts files | Field names, data shapes |
| API | Read client/api code | URL construction, headers |
| Backend | Read handler code | Parameters, service calls |
| Business Logic | Read service code | Query logic, validation |
| Database | Read migrations | Table structure, constraints |

---

## **Common Bug Patterns**

1. **Field name mismatches** → `destination_url` vs `destinationUrl`
2. **Missing authorization** → Endpoint protected but no auth check
3. **URL construction errors** → Wrong base URL, wrong parameters
4. **Type mismatches** → String vs UUID, wrong array type
5. **Constraint conflicts** → Unique constraint blocks valid data
6. **Scope issues** → Global vs per-user uniqueness
7. **Parameter passing** → Function signature changes break callers
8. **Schema/code drift** → Migration doesn't match service expectations

---

## **The 4-Step Quick Process**

When debugging an issue:

1. **Trace the feature** (user action → database)
2. **Examine each layer** for obvious issues
3. **Form a hypothesis** based on evidence
4. **Fix minimally & check ripple effects**

This approach catches most bugs within 5-10 minutes because bugs usually appear at layer boundaries (frontend ↔ backend, code ↔ database).

---

## **Real Examples from LinkPulse**

### **Bug 1: Form Inputs Not Visible**
- **Layer:** UI/Component
- **Evidence:** `border-gray-300` class without `border` class
- **Hypothesis:** Missing Tailwind utility for border
- **Fix:** Added `border` class to inputs
- **Time:** ~2 minutes

### **Bug 2: Can't Create Links**
- **Layer:** Database
- **Evidence:** Created migration with `UNIQUE(slug)`
- **Hypothesis:** Global slug constraint prevents multiple users from having same slug
- **Fix:** Changed to composite `UNIQUE(user_id, slug)` constraint
- **Time:** ~5 minutes

### **Bug 3: Redirect Logs Out User**
- **Layer:** Frontend ↔ Backend mismatch
- **Evidence:** Frontend uses `user_id` but backend route expects `username`
- **Hypothesis:** URL parameter mismatch
- **Fix:** Updated handler to extract username, service to use userID
- **Time:** ~10 minutes

### **Bug 4: Click Link Redirects to Login**
- **Layer:** URL construction
- **Evidence:** Short URL uses `window.location.origin` (port 5173 frontend)
- **Hypothesis:** Backend API on port 8080 is being ignored
- **Fix:** Changed to hardcoded backend URL
- **Time:** ~3 minutes

---

**Total bugs fixed: 4 | Total time: ~20 minutes | Success rate: 100%**
