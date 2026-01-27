include(FetchContent)

function(codegen_init GIT_TAG)
    FetchContent_Declare(
        codegen
        GIT_REPOSITORY https://github.com/schraf/codegen.git
        GIT_TAG        ${GIT_TAG}
    )
    FetchContent_MakeAvailable(codegen)
    
    set(CODEGEN_BIN "${CMAKE_BINARY_DIR}/bin/codegen${CMAKE_EXECUTABLE_SUFFIX}" CACHE INTERNAL "")
    
    add_custom_command(
        OUTPUT "${CODEGEN_BIN}"
        COMMAND go build -o "${CODEGEN_BIN}" ./cmd/...
        WORKING_DIRECTORY "${codegen_SOURCE_DIR}"
        COMMENT "Building code generation tool..."
        VERBATIM
    )
    
    # Create a target for the executable so other steps can depend on it
    add_custom_target(build_codegen DEPENDS "${CODEGEN_BIN}")
endfunction()

function(codegen_project TARGET_NAME)
    set(options "")
    set(oneValueArgs PROJECT_FILE)
    set(multiValueArgs OUTPUTS DEPENDS)

    cmake_parse_arguments(ARG "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})

    add_custom_command(
        OUTPUT ${ARG_OUTPUTS}
        COMMAND "${CODEGEN_BIN}" -project "${ARG_PROJECT_FILE}"
        DEPENDS "${CODEGEN_BIN}" "${ARG_PROJECT_FILE}" ${ARG_DEPENDS}
        WORKING_DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}"
        COMMENT "Generating code for ${TARGET_NAME} from ${ARG_PROJECT_FILE}"
        VERBATIM
    )
                                                            
    target_sources(${TARGET_NAME} PRIVATE ${ARG_OUTPUTS})
    add_dependencies(${TARGET_NAME} build_codegen)
endfunction()

