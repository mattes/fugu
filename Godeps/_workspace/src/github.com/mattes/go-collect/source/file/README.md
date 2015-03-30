# File source

Currently supports YAML files only. 
[More example files.](https://github.com/mattes/fugu/tree/v1/examples)

## Example URLs

```
--source=file:///absolute/path/to/file.yml
--source=file://relative/path/to/file.yml
--source=file://../relative/path/to/file.yml
--source=file://file.yml
```

## Default label

* If you don't ask for a specific label, it will
  * return label ``default`` if found
  * or return first found label

* If you ask for a specific label, it will
  * return this specific label if found
  * and just don't return anything else if not found