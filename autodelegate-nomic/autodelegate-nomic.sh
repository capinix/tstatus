#!/bin/bash -e

# requires awk bc

# Default values changes can be made in the sourced conf file
MINIMUM_AMOUNT=1000000
MINIMUM_BALANCE=1000000

BINARY=/root/.cargo/bin/nomic

THIS_FILE=${BASH_SOURCE[0]}
CONFIG_FILE="${THIS_FILE%.*}.conf"

# Sourced conf file for changes to default parameters
[ -f ${CONFIG_FILE} ] && source ${CONFIG_FILE}

# Claim rewards, If the claim gives an error exit
${BINARY} claim || exit

# Lets see how much we have after claiming our rewards
BALANCE=$(${BINARY} balance | awk '/balance/{print $2}')

if (( ${BALANCE} < ( ${MINIMUM_AMOUNT} + ${MINIMUM_BALANCE} ) )) ; then

  # get the current delegations to the validator
  DELEGATED=$(${BINARY} delegations | awk -F' |=' "/${VALIDATOR}/"'{print $4}')

  # convert uNOM to NOM with 2 decimal places
  DELEGATED=$(echo "scale=2; ${DELEGATED} / 1000000" | bc)

  # Not enough to delegate
  env LC_ALL=en_US.UTF-8 printf "\u274c \e[31m%'.0f uNOM\e[0m not enough to delegate. Self delegation \e[3m%'.2f NOM\e[0m\n" ${BALANCE} ${DELEGATED}

else

  # we have enough to delegate, calculate the quantity to delegate
  # total in wallet less the minimum balance we like to leave
  # in uNOM
  QUANTITY=$(( ${BALANCE} - ${MINIMUM_BALANCE} ))

  # delegate to our validator
  ${BINARY} delegate ${VALIDATOR} ${QUANTITY}

  # convert uNOM to NOM with 2 decimal places
  QUANTITY=$(echo "scale=2; ${QUANTITY} / 1000000" | bc)

  # get the current delegations to the validator
  DELEGATED=$(${BINARY} delegations | awk -F' |=' "/${VALIDATOR}/"'{print $4}')

  # convert uNOM to NOM with 2 decimal places
  DELEGATED=$(echo "scale=2; ${DELEGATED} / 1000000" | bc)

  # log success
  env LC_ALL=en_US.UTF-8 printf "\u2705 Delegated \e[32m%'.2f NOM\e[0m. Self delegation \e[3m%'.2f NOM\e[0m\n" ${QUANTITY} ${DELEGATED}

fi
