
data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "../src/infrastructure/persistence/models",
    "--dialect", "postgres",
  ]
}

env "gorm-dev" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/17/dev"
  migration {
    dir = "file://versions/dev/"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}


env "gorm-prod" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/17/dev"
  migration {
    dir = "file://versions/prod/"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

