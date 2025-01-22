#!/bin/bash

 PLI_VARS_FILE='pluginGoVariable.vars'
cp sample.vars ~/$PLI_VARS_FILE
chmod 600 ~/$PLI_VARS_FILE

source ~/$PLI_VARS_FILE
