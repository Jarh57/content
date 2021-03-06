package content

import (
	"testing"

	"github.com/stretchr/testify/require"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

var (
	translatedProviderSess   sqlbuilder.Database
	translatedProviderModels db.Collection
)

type testTranslatedProviderModel struct {
	ID          int64              `db:"id,omitempty"`
	Name        TranslatedProvider `db:"name"`
	Description TranslatedProvider `db:"description"`
}

func initTranslatedProviderDB(t *testing.T) {
	cnf := &mysql.ConnectionURL{
		User:     "dev-user",
		Password: "dev-password",
		Host:     "localhost",
		Database: "test",
		Options: map[string]string{
			"charset":   "utf8mb4",
			"collation": "utf8mb4_bin",
			"parseTime": "true",
		},
	}
	var err error
	translatedProviderSess, err = mysql.Open(cnf)
	require.Nil(t, err)

	_, err = translatedProviderSess.Exec(`DROP TABLE IF EXISTS translated_test`)
	require.Nil(t, err)

	_, err = translatedProviderSess.Exec(`
    CREATE TABLE translated_test (
      id INT(11) NOT NULL AUTO_INCREMENT,
      name JSON,
      description JSON,

      PRIMARY KEY(id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
  `)
	require.Nil(t, err)

	translatedProviderModels = translatedProviderSess.Collection("translated_test")

	require.Nil(t, translatedProviderModels.Truncate())
}

func finishTranslatedProviderDB() {
	translatedProviderSess.Close()
}

func TestLoadSaveTranslatedProvider(t *testing.T) {
	initTranslatedProviderDB(t)
	defer finishTranslatedProviderDB()

	model := new(testTranslatedProviderModel)
	require.Nil(t, translatedProviderModels.InsertReturning(model))

	require.EqualValues(t, model.ID, 1)

	other := new(testTranslatedProviderModel)
	require.Nil(t, translatedProviderModels.Find(1).One(other))
}

func TestLoadSaveTranslatedProviderWithContent(t *testing.T) {
	initTranslatedProviderDB(t)
	defer finishTranslatedProviderDB()

	model := &testTranslatedProviderModel{
		Name: TranslatedProvider{
			"altipla": Translated{
				"es": "foo-es",
				"en": "foo-en",
			},
			"hotelbeds": Translated{
				"es": "bar-es",
				"en": "bar-en",
			},
		},
	}
	require.Nil(t, translatedProviderModels.InsertReturning(model))

	other := new(testTranslatedProviderModel)
	require.Nil(t, translatedProviderModels.Find(1).One(other))

	require.Equal(t, other.Name["altipla"]["es"], "foo-es")
	require.Equal(t, other.Name["altipla"]["en"], "foo-en")
	require.Equal(t, other.Name["hotelbeds"]["es"], "bar-es")
	require.Equal(t, other.Name["hotelbeds"]["en"], "bar-en")
}

func TestTranslatedProviderGlobalChain(t *testing.T) {
	SetGlobalProviderChain([]string{"altipla", "hotelbeds", "dingus"})

	tests := []struct {
		content TranslatedProvider
		chain   Translated
	}{
		{
			TranslatedProvider{
				"altipla": Translated{
					"es": "foo-es",
					"en": "foo-en",
				},
				"hotelbeds": Translated{
					"es": "bar-es",
					"en": "bar-en",
				},
				"dingus": Translated{
					"es": "baz-es",
					"en": "baz-en",
				},
			},
			Translated{
				"es": "foo-es",
				"en": "foo-en",
			},
		},
		{
			TranslatedProvider{
				"altipla": Translated{
					"es": "foo-es",
					"en": "foo-en",
				},
				"hotelbeds": Translated{
					"es": "bar-es",
					"en": "bar-en",
				},
				"dingus": Translated{
					"es": "baz-es",
					"en": "baz-en",
					"it": "baz-it",
				},
			},
			Translated{
				"es": "foo-es",
				"en": "foo-en",
				"it": "baz-it",
			},
		},
		{
			TranslatedProvider{
				"altipla": Translated{
					"es": "foo-es",
				},
				"hotelbeds": Translated{
					"es": "bar-es",
					"en": "bar-en",
				},
				"dingus": Translated{
					"es": "baz-es",
					"en": "baz-en",
				},
			},
			Translated{
				"es": "foo-es",
				"en": "bar-en",
			},
		},
		{
			TranslatedProvider{
				"altipla": Translated{
					"es": "foo-es",
				},
				"hotelbeds": Translated{
					"es": "bar-es",
					"en": "bar-en",
				},
				"tripadvisor": Translated{
					"it": "qux-it",
				},
			},
			Translated{
				"es": "foo-es",
				"en": "bar-en",
			},
		},
		{
			TranslatedProvider{
				"tripadvisor": Translated{
					"it": "qux-it",
				},
			},
			Translated{},
		},
	}
	for _, test := range tests {
		require.Equal(t, test.content.Chain(), test.chain)
	}
}

func TestTranslatedProviderCustomChain(t *testing.T) {
	SetGlobalProviderChain([]string{"dingus", "tirpadvisor", "other"})

	content := TranslatedProvider{
		"altipla": Translated{
			"es": "foo-es",
			"en": "foo-en",
		},
		"hotelbeds": Translated{
			"es": "bar-es",
			"en": "bar-en",
		},
		"dingus": Translated{
			"es": "baz-es",
			"en": "baz-en",
		},
	}
	require.Equal(t, content.CustomChain([]string{"altipla", "hotelbeds", "dingus"}), Translated{
		"es": "foo-es",
		"en": "foo-en",
	})
}

func TestTranslatedProviderSetValue(t *testing.T) {
	content := make(TranslatedProvider)
	content.SetValue("hotelbeds", "es", "value")
	require.Equal(t, content["hotelbeds"]["es"], "value")

	content.SetValue("hotelbeds", "es", "value2")
	require.Equal(t, content["hotelbeds"]["es"], "value2")
}

func TestTranslatedProviderGetProvider(t *testing.T) {
	content := TranslatedProvider{
		"altipla": Translated{"es": "foo"},
	}
	require.Equal(t, content.Provider("hotelbeds"), Translated{})

	require.Equal(t, content.Provider("altipla"), Translated{"es": "foo"})
}
