package httpapi

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sourcegraph/log"
	"github.com/sourcegraph/sourcegraph/internal/cody"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/ssc"
)

// newSSCRefreshCodyRateLimitHandler returns an http.Handler to trigger cody's rate limit refresh for a user
// TODO(sourcegraph#59625) remove as part of adding SAMSActor source
func newSSCRefreshCodyRateLimitHandler(logger log.Logger, db database.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		samsAccountID := mux.Vars(r)["samsAccountID"]
		if samsAccountID == "" {
			http.Error(w, "missing uuid", http.StatusBadRequest)
			return
		}

		oidcAccounts, err := db.UserExternalAccounts().List(ctx, database.ExternalAccountsListOptions{
			AccountID:   samsAccountID,
			ServiceType: "openidconnect",
			ServiceID:   fmt.Sprintf("https://%s", ssc.GetSAMSHostName()),
			LimitOffset: &database.LimitOffset{
				Limit: 1,
			},
		})
		if err != nil {
			logger.Error("error getting oidc accounts", log.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(oidcAccounts) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		userID := oidcAccounts[0].UserID

		if err := cody.RefreshGatewayRateLimits(ctx, userID, db); err != nil {
			logger.Error("error refreshing gateway rate limits", log.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})
}
