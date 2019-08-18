# srcreader

Responsibilities:

* Clone a GitHub repository.
* Access each _.go_ file.
* Create an AST.
* Process each AST to extract required information (pre-processing).
* Process each AST to:
  * Extract tokens (identifiers)
  * Split those tokens
  * Expand those splittings
* Generate a newly and modified AST.
* Register metrics.
